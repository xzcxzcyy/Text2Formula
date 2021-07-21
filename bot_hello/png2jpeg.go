package main

import (
    "image"
    "image/color"
    "image/draw"
    "image/jpeg"
    "image/png"
    "log"
    "os"
)

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

    // create a new Image with the same dimension of PNG image
    newImg := image.NewRGBA(imgSrc.Bounds())

    // we will use white background to replace PNG's transparent background
    // you can change it to whichever color you want with
    // a new color.RGBA{} and use image.NewUniform(color.RGBA{<fill in color>}) function

    draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

    // paste PNG image OVER to newImage
    draw.Draw(newImg, newImg.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Over)

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
