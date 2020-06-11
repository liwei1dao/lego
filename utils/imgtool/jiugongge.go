package imgtool

import (
	"image"
	"image/color"
	"io"

	"github.com/disintegration/imaging"
)

type Point struct {
	x, y int
}

func BuildJiuGong(src []io.Reader, dst io.ReadWriter, format imaging.Format, opts ...imaging.EncodeOption) (err error) {
	imagePoints := getXy(len(src))
	width := getWidth(len(src))

	background := imaging.New(132, 132, color.RGBA{233, 233, 233, 255})
	for i, v := range imagePoints {
		x := v.x
		y := v.y
		if src, err := imaging.Decode(src[i]); err == nil {
			src = imaging.Resize(src, width, width, imaging.Lanczos)
			background = imaging.Paste(background, src, image.Pt(x, y))
		}
	}
	err = imaging.Encode(dst, background, format, opts...)
	return nil
}

func BuildJiuGongToFlie(dstpath string, src []io.Reader, opts ...imaging.EncodeOption) (err error) {
	imagePoints := getXy(len(src))
	width := getWidth(len(src))

	background := imaging.New(132, 132, color.RGBA{233, 233, 233, 255})
	for i, v := range imagePoints {
		x := v.x
		y := v.y
		if src, err := imaging.Decode(src[i]); err == nil {
			src = imaging.Resize(src, width, width, imaging.Lanczos)
			background = imaging.Paste(background, src, image.Pt(x, y))
		}
	}
	err = imaging.Save(background, dstpath, opts...)
	return nil
}

func getXy(size int) []*Point {
	s := make([]*Point, size)
	var _x, _y int

	if size == 1 {
		_x, _y = 6, 6
		s[0] = &Point{_x, _y}
	}
	if size == 2 {
		_x, _y = 4, 4
		s[0] = &Point{_x, 132/2 - 60/2}
		s[1] = &Point{60 + 2*_x, 132/2 - 60/2}
	}
	if size == 3 {
		_x, _y = 4, 4
		s[0] = &Point{132/2 - 60/2, _y}
		s[1] = &Point{_x, 60 + 2*_y}
		s[2] = &Point{60 + 2*_y, 60 + 2*_y}
	}
	if size == 4 {
		_x, _y = 4, 4
		s[0] = &Point{_x, _y}
		s[1] = &Point{_x*2 + 60, _y}
		s[2] = &Point{_x, 60 + 2*_y}
		s[3] = &Point{60 + 2*_y, 60 + 2*_y}
	}
	if size == 5 {
		_x, _y = 3, 3
		s[0] = &Point{(132 - 40*2 - _x) / 2, (132 - 40*2 - _y) / 2}
		s[1] = &Point{((132-40*2-_x)/2 + 40 + _x), (132 - 40*2 - _y) / 2}
		s[2] = &Point{_x, ((132-40*2-_x)/2 + 40 + _y)}
		s[3] = &Point{(_x*2 + 40), ((132-40*2-_x)/2 + 40 + _y)}
		s[4] = &Point{(_x*3 + 40*2), ((132-40*2-_x)/2 + 40 + _y)}
	}
	if size == 6 {
		_x, _y = 3, 3
		s[0] = &Point{_x, ((132 - 40*2 - _x) / 2)}
		s[1] = &Point{(_x*2 + 40), ((132 - 40*2 - _x) / 2)}
		s[2] = &Point{(_x*3 + 40*2), ((132 - 40*2 - _x) / 2)}
		s[3] = &Point{_x, ((132-40*2-_x)/2 + 40 + _y)}
		s[4] = &Point{(_x*2 + 40), ((132-40*2-_x)/2 + 40 + _y)}
		s[5] = &Point{(_x*3 + 40*2), ((132-40*2-_x)/2 + 40 + _y)}
	}

	if size == 7 {
		_x, _y = 3, 3
		s[0] = &Point{(132 - 40) / 2, _y}
		s[1] = &Point{_x, (_y*2 + 40)}
		s[2] = &Point{(_x*2 + 40), (_y*2 + 40)}
		s[3] = &Point{(_x*3 + 40*2), (_y*2 + 40)}
		s[4] = &Point{_x, (_y*3 + 40*2)}
		s[5] = &Point{(_x*2 + 40), (_y*3 + 40*2)}
		s[6] = &Point{(_x*3 + 40*2), (_y*3 + 40*2)}
	}
	if size == 8 {
		_x, _y = 3, 3
		s[0] = &Point{(132 - 80 - _x) / 2, _y}
		s[1] = &Point{((132-80-_x)/2 + _x + 40), _y}
		s[2] = &Point{_x, (_y*2 + 40)}
		s[3] = &Point{(_x*2 + 40), (_y*2 + 40)}
		s[4] = &Point{(_x*3 + 40*2), (_y*2 + 40)}
		s[5] = &Point{_x, (_y*3 + 40*2)}
		s[6] = &Point{(_x*2 + 40), (_y*3 + 40*2)}
		s[7] = &Point{(_x*3 + 40*2), (_y*3 + 40*2)}
	}
	if size == 9 {
		_x, _y = 3, 3
		s[0] = &Point{_x, _y}
		s[1] = &Point{_x*2 + 40, _y}
		s[2] = &Point{_x*3 + 40*2, _y}
		s[3] = &Point{_x, (_y*2 + 40)}
		s[4] = &Point{(_x*2 + 40), (_y*2 + 40)}
		s[5] = &Point{(_x*3 + 40*2), (_y*2 + 40)}
		s[6] = &Point{_x, (_y*3 + 40*2)}
		s[7] = &Point{(_x*2 + 40), (_y*3 + 40*2)}
		s[8] = &Point{(_x*3 + 40*2), (_y*3 + 40*2)}
	}
	return s
}

func getWidth(size int) int {
	var width int
	if size == 1 {
		width = 120
	}
	if size > 1 && size <= 4 {
		width = 60
	}
	if size >= 5 {
		width = 40
	}
	return width
}
