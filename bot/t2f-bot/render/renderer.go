package render

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/jpeg"
    "image/png"
    "log"
    "os"
    "os/exec"
)

const (
    minImageWidth  int = 720
    minImageHeight int = 360
)

type Renderer struct {
    SvgDir     string
    PngDir     string
    JpgDir     string
    MathjaxDir string
}

func max(a, b int) int {
    if a > b {
        return a
    } else {
        return b
    }
}

func png2Jpg(pngFilePath string, jpgFilePath string) (sizeInfo image.Point, retErr error) {
    pngImgFile, err := os.Open(pngFilePath)

    if err != nil {
        return image.Point{}, err
    }

    defer pngImgFile.Close()

    // create image from PNG file
    imgSrc, err := png.Decode(pngImgFile)

    if err != nil {
        return image.Point{}, err
    }

    srcImgSize := imgSrc.Bounds().Size()
    newImgBounds := image.Rectangle{
        Min: image.Point{},
        Max: image.Point{
            X: max(srcImgSize.X, minImageWidth) + 5,
            Y: max(srcImgSize.Y, minImageHeight) + 5,
        },
    }
    newImg := image.NewRGBA(newImgBounds)
    draw.Draw(newImg, newImgBounds, &image.Uniform{C: color.White}, image.Point{}, draw.Src)
    dp := image.Point{
        X: newImgBounds.Min.X + (newImgBounds.Size().X-srcImgSize.X)/2,
        Y: newImgBounds.Min.Y + (newImgBounds.Size().Y-srcImgSize.Y)/2,
    }
    drawBounds := image.Rectangle{Min: dp, Max: dp.Add(srcImgSize)}
    draw.Draw(newImg, drawBounds, imgSrc, imgSrc.Bounds().Min, draw.Over)

    // create new out JPEG file
    jpgImgFile, err := os.Create(jpgFilePath)

    if err != nil {
        log.Println("Cannot create JPEG file.")
        return image.Point{}, err
    }

    defer jpgImgFile.Close()

    var opt jpeg.Options
    opt.Quality = 80

    // convert newImage to JPEG encoded byte and save to jpgImgFile
    // with quality = 80
    err = jpeg.Encode(jpgImgFile, newImg, &opt)

    //err = jpeg.Encode(jpgImgFile, newImg, nil) -- use nil if ignore quality options

    if err != nil {
        return image.Point{}, err
    }

    //log.Println("Converted PNG file to JPEG file")
    return newImg.Bounds().Size(), nil
}

func (r *Renderer) RenderTex(queryID string, formula string) (filePath string, sizeInfo image.Point, retErr error) {
    curSvgFilePath := fmt.Sprintf("%s/%s.svg", r.SvgDir, queryID)
    curPngFilePath := fmt.Sprintf("%s/%s.png", r.PngDir, queryID)
    curJpgFilePath := fmt.Sprintf("%s/%s.jpg", r.JpgDir, queryID)
    perfTex2Svg := exec.Command(r.MathjaxDir+"/tex2svg", formula, curSvgFilePath)
    //log.Println(perfTex2Svg.Args)
    perfTex2Svg.Dir = r.MathjaxDir
    err := perfTex2Svg.Run()
    if err != nil {
        log.Printf("during perfTex2Svg: %v", err)
        return "", image.Point{}, err
    }
    //log.Println(perfTex2Svg.Args)
    perfSvg2Png := exec.Command("cairosvg", curSvgFilePath, "-o", curPngFilePath, "-s", "2.5")
    //log.Println(perfSvg2Png.Args)
    err = perfSvg2Png.Run()
    if err != nil {
        log.Printf("during perfSvg2Png: %v", err)
        return "", image.Point{}, err
    }
    imgSize, err := png2Jpg(curPngFilePath, curJpgFilePath)
    if err != nil {
        log.Printf("during png2Jpg: %v", err)
        return "", image.Point{}, err
    }
    return curJpgFilePath, imgSize, nil
}
