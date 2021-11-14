package raytrace

import(
	"fmt"
	"image/color"
	"math"
)

type Camera struct{
	Height, Width int
	Origin Vec3
	Lower_left_corner Vec3
	Horizontal Vec3
	Vertical Vec3
}

func NewCamera(lookfrom, lookat Vec3, width int) Camera {

	//Image
	vfov := 90.0
	aspect_ratio := 16.0 / 9.0
	height := int(float64(width) / aspect_ratio)
	fmt.Printf("image res", width, height)
	
	// Camera
	cam := Camera{}
	cam.Width = width
	cam.Height = height
	theta := float64(Deg_to_Rad(vfov))
	h := math.Tan(theta/2) // half height
	_ = h
	viewport_height := 2.0 *h
	viewport_width := aspect_ratio * viewport_height
	
	vup := NewVec3(0,1,0)
	w := (lookfrom.Subtr(lookat)).UnitVec()
	u := (w.Cross(vup)).UnitVec()
	v := u.Cross(w)
	// fmt.Println("w,u,v", w, u, v, viewport_width, viewport_height)
	cam.Origin = lookfrom

	// NOTE: Start from u, or v vector and multiple by viewport_width float
	//       In order to scale each axis - this was the issue in the previous commit
	cam.Horizontal = u.MultF(float32(viewport_width))
	cam.Vertical = v.MultF(float32(viewport_height))
	// fmt.Println("hor/ver", cam.Horizontal, cam.Vertical)
	
	o := cam.Origin.Subtr(cam.Horizontal.DivF(2.0))
	o = o.Subtr(cam.Vertical.DivF(2.0))
	cam.Lower_left_corner = o.Subtr(w)
	// fmt.Println("llc:", cam.Lower_left_corner)
	return cam
}

func (c Camera) GetRay(u,v float32) Ray {
	u_horiz := c.Horizontal.MultF(u)
	v_vert := c.Vertical.MultF(v)
	dir := c.Lower_left_corner.Add(u_horiz)
	dir = dir.Add(v_vert)
	dir = dir.Subtr(c.Origin)
	return NewRay(c.Origin, dir)
}

func Write_color(cd Vec3, samples int) color.RGBA {
	scale := float32(1.0) / float32(samples)
	R := cd.At(0) * scale
	G := cd.At(1) * scale
	B := cd.At(2) * scale
	R = Clamp(R, float32(0.0), float32(0.9999))
	G = Clamp(G, float32(0.0), float32(0.9999))
	B = Clamp(B, float32(0.0), float32(0.9999))
	return color.RGBA{uint8(R*255), uint8(G*255), uint8(B*255), 255}	
}

func RayColorArray(r *Ray, world HittableList) Vec3 {
	rec := HitRecord{NewVec3(0,0,0), NewVec3(0,0,0), 1.0, true, -1}
	hit := world.Hit(r, 0, float32(math.Inf(1.0)), &rec)
	
	if hit {
		N := (rec.Normal.Add(NewVec3(1,1,1))).MultF(float32(0.5))
		// N = N.UnitVec()
		return N
	}
	// Background
	unit_direction := r.Direction().UnitVec()
	t := float32(0.5 * (unit_direction.At(1) + 1.0))
	sky := NewVec3(0.5, 0.7, 1.0).MultF(t)
	sky = sky.Add(NewVec3(1, 1, 1).MultF(1 - t))
	return sky
}

func RayColorBVH(r *Ray, world *BVH_node) Vec3 {
	rec := HitRecord{NewVec3(0,0,0), NewVec3(0,0,0), 1.0, true, -1}
	hit := world.Hit(r, 0, float32(math.Inf(1.0)), &rec)
	
	if hit {
		N := (rec.Normal.Add(NewVec3(1,1,1))).MultF(float32(0.5))
		// N = N.UnitVec()
		return N
	}
	// Background
	unit_direction := r.Direction().UnitVec()
	t := float32(0.5 * (unit_direction.At(1) + 1.0))
	sky := NewVec3(0.5, 0.7, 1.0).MultF(t)
	sky = sky.Add(NewVec3(1, 1, 1).MultF(1 - t))
	return sky
}
