package main

import (
    "./resize"
    "image"
    "os"
    "io"
    "log"
    "math"
    "strconv"
    "image/png"
    "image/draw"
    "net/http"
    "regexp"
    "path"
    "path/filepath"
)

func Resize(writer io.Writer, path string, w int, h int) (error) {
    file, err := os.Open(path)

    if err != nil {
        return err
    }

    defer file.Close()

    picture, _, err := image.Decode(file)

    if err != nil {
        return err
    }

    bounds := picture.Bounds()

    size := bounds.Size()

    iw, ih := size.X, size.Y

    if iw > w {
        ih = int(math.Max(float64(size.Y) * float64(w) / float64(size.X), 1.0))
        iw = w
    }
    if ih > h {
        iw = int(math.Max(float64(size.X) * float64(h) / float64(size.Y), 1.0))
        ih = h
    }

    resized := resize.Resize(uint(iw), uint(ih), picture, resize.Bilinear)

    output := image.NewRGBA(image.Rect(0, 0, w, h))

    draw.Draw(output, image.Rect(0, 0, iw, ih).Add(image.Pt((w - iw) / 2, (h - ih) / 2)), resized, image.ZP, draw.Over)

    png.Encode(writer, output)

    return nil
}

func main() {
    imageDirectory := "./images"
    basepath, err := filepath.Abs(imageDirectory)

    if (err != nil) {
        log.Fatal(err)
    }

    queryPattern := regexp.MustCompile(`^/(?P<width>\d+)/(?P<height>\d+)/(?P<path>.+(?:png|jpe?g))$`)

    server := http.Server {
        Addr: ":8080",
        Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // log.Printf(`%q`, html.EscapeString(r.URL.Path))

            match := queryPattern.FindStringSubmatch(r.URL.Path)

            if match != nil {
                width, _ := strconv.Atoi(match[1])
                height, _ := strconv.Atoi(match[2])

                filename := path.Join(basepath, filepath.Clean(match[3]))

                // do not expose the absolute path of the file in error messages

                filename, err = filepath.Rel(basepath, filename)

                if err == nil {
                    err = Resize(w, path.Join(imageDirectory, filename), width, height)
                }

                if err != nil {
                    http.Error(w, err.Error(), 500)
                }
            } else {
                http.NotFound(w, r)
            }
        }),
    }

    log.Printf("Listening on %s...", server.Addr)
    server.ListenAndServe()
}