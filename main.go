package main

import (
	 "github.com/apaladiychuk/store"
	"os"
	"io/ioutil"
	"log"
	"fmt"
	"github.com/astaxie/beego/logs"
	"image/jpeg"
	"github.com/nfnt/resize"
	"sync"
	"time"
	"strconv"
)
var (
	sourceFolder, destFolder string
	maxWidth int
)

func init() {
	store.Init(store.StoreInLocalDirectory)
}

func main() {
	c := make( chan string )

	if len(os.Args) != 4 {
		fmt.Println("Wrong argument ")
		fmt.Printf("usage  %s <source folder> <destination folder> <max width> \n ", os.Args[0])
		return
	}
	sourceFolder = os.Args[1]
	destFolder = os.Args[2]
	var err error
	if maxWidth ,err  = strconv.Atoi( os.Args[3] ); err != nil   {
		fmt.Println(" error parameter  : ",  err.Error())
		return
	}
	if maxWidth == 0 {
		fmt.Println(" error parameter  : max width connot be 0 " )
		return
	}
	fmt.Println ( sourceFolder)
	fmt.Println ( destFolder)
	files, err := ioutil.ReadDir(sourceFolder)
	if err != nil {
		log.Fatal(err)
	}
	go logger( &c )
	st := time.Now()
	wg := sync.WaitGroup{}

	wg.Add(len(files) )
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		//
		 processImage(&wg, &c, f.Name())
	}

	fmt.Println(" leave loop  ")
	wg.Wait()

     duration :=time.Since(st)
	//float32(st.Nanosecond() - fn.Nanosecond() ) / 1000000000
	fmt.Printf(" after wait  %f  ", duration.Seconds()  )
	c <- "EXIT"

}

func processImage( wg *sync.WaitGroup , c *chan string  , fileName string  ){
	file, err := os.Open(sourceFolder + fileName )
	defer func() {
		wg.Done()
	}()

	if err != nil {
		logs.Error(err.Error())
		return
	}
	defer file.Close()
	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	size := img.Bounds().Size()
	koef := float32( size.X ) / float32(maxWidth )
	msg := fmt.Sprintf( " file = %s , width = %d , height = %d  koef = %.3f \n ", fileName, size.X, size.Y , koef )
	*c <- msg
	if koef < 1  {
		koef = 1
	}

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	//m := resize.Resize(1000, 0, img, resize.Lanczos3)

	m := resize.Thumbnail(uint( float32(size.X) / koef )  , uint( float32(size.Y ) / koef ) , img, resize.Lanczos3)
	out, err := os.Create(destFolder + fileName )
	if err != nil {
		logs.Error(err.Error())
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

}
func logger( c *chan string ){
	fmt.Println(" start logger  ")
	for true {
		message := <-  *c
		if message == "EXIT" {
			fmt.Println(" finish application ")
			break
		} else {
			fmt.Print(message)
		}

	}
}