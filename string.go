package memdraw

/*
func String(dst draw.Image, p image.Point, src image.Image, sp image.Point, ft *font.Font, s []byte) int {
	for _, b := range s {
		mask := ft.Char(b)
		if mask == nil {
			panic("StringBG")
		}
		r := mask.Bounds()
		draw.DrawMask(dst, r.Add(p), src, sp, mask, mask.Bounds().Min, draw.Over)
		p.X += r.Dx() + ft.stride
	}
	return p.X
}
*/
