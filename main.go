package main

import (
	"gocv.io/x/gocv"
	"gonum.org/v1/gonum/mat"
	"image"
	"math"
)

func main() {
	original_window := gocv.NewWindow("Original")
	original_image := gocv.IMRead("sample/original.png", gocv.IMReadColor)
	for {
		original_window.IMShow(original_image)
		if original_window.WaitKey(1) >= 0 {
			break
		}
	}

	pts := mat.NewDense(4, 2, []float64{
		69, 145,
		97, 735,
		938, 723,
		971, 170,
	})

	tranformed_window := gocv.NewWindow("Transformed")
	transformed_image := FourPointTransform(original_image, pts)
	for {
		tranformed_window.IMShow(transformed_image)
		if tranformed_window.WaitKey(1) >= 0 {
			break
		}
	}
}
func FourPointTransform(img gocv.Mat, pts *mat.Dense) gocv.Mat{
	rect := orderPoints(pts)
	tl := rect.RawRowView(0)
	tr := rect.RawRowView(1)
	br := rect.RawRowView(2)
	bl := rect.RawRowView(3)

	// compute the width of the new image, which will be the
	// maximum distance between bottom-right and bottom-left
	// x-coordiates or the top-right and top-left x-coordinates
	widthA := math.Sqrt(math.Pow((br[0] - bl[0]), 2) + math.Pow((br[1] - bl[1]), 2))
	widthB := math.Sqrt(math.Pow((tr[0] - tl[0]), 2) + math.Pow((tr[1] - tl[1]), 2))
	maxWidth := int(math.Max(widthA, widthB))

	// compute the height of the new image, which will be the
	// maximum distance between the top-right and bottom-right
	// y-coordinates or the top-left and bottom-left y-coordinates
	heightA := math.Sqrt(math.Pow((tr[0] - br[0]),2) + math.Pow((tr[1] - br[1]), 2))
	heightB := math.Sqrt(math.Pow((tl[0] - bl[0]), 2) + math.Pow((tl[1] - bl[1]), 2))
	maxHeight := int(math.Max(heightA, heightB))

	// now that we have the dimensions of the new image, construct
	// the set of destination points to obtain a "birds eye view",
	// (i.e. top-down view) of the image, again specifying points
	// in the top-left, top-right, bottom-right, and bottom-left
	// order

	//dst = np.array([
	//	[0, 0],
	//[maxWidth - 1, 0],
	//[maxWidth - 1, maxHeight - 1],
	//[0, maxHeight - 1]], dtype="float32")
	//M = cv2.getPerspectiveTransform(rect, dst)
	//warped = cv2.warpPerspective(image, M, (maxWidth, maxHeight))
	dst := mat.NewDense(4, 2, []float64{
		0, 0,
		(float64(maxWidth) - 1), 0,
		(float64(maxWidth) - 1), (float64(maxHeight) - 1),
		0, (float64(maxHeight) - 1),
	})

	M := gocv.GetPerspectiveTransform(convertDenseToImagePoint(rect), convertDenseToImagePoint(dst))
	gocv.WarpPerspective(img, &img, M, image.Point{X: maxWidth, Y: maxHeight})

	return img
}

func orderPoints(pts *mat.Dense) *mat.Dense{
	// initialzie a list of coordinates that will be ordered
	// such that the first entry in the list is the top-left,
	// the second entry is the top-right, the third is the
	// bottom-right, and the fourth is the bottom-left

	rect := mat.NewDense(4, 2, nil)

	// the top-left point will have the smallest sum, whereas
	// the bottom-right point will have the largest sum
	sumMinIndex, sumMaxIndex := findMinMaxSumIndex(*pts)
	rect.SetRow(0, pts.RawRowView(sumMinIndex))
	rect.SetRow(2, pts.RawRowView(sumMaxIndex))

	// now, compute the difference between the points, the
	// top-right point will have the smallest difference,
	// whereas the bottom-left will have the largest difference
	diffMinIndex, diffMaxIndex := findMinMaxDiffIndex(*pts)
	rect.SetRow(1, pts.RawRowView(diffMinIndex))
	rect.SetRow(3, pts.RawRowView(diffMaxIndex))

	// return the ordered coordinates
	return rect
}

func findMinMaxSumIndex(pts mat.Dense) (int, int){
	r, c := pts.Dims()

	maxIndex := 0
	maxValue := 0.0
	minIndex := 0
	minValue := 0.0

	for i := 0; i < r; i++ {
		row := pts.RowView(i)
		sum := 0.0
		for j := 0; j < c; j++ {
			sum += row.AtVec(j)
		}

		if (i == 0 ) {
			maxValue = sum
			minValue = sum
		}

		//Find max value and index
		if (sum > maxValue) {
			maxValue = sum
			maxIndex = i
		}
		//Find min value and index
		if (sum < minValue) {
			minValue = sum
			minIndex = i
		}
	}
	return minIndex, maxIndex
}

func findMinMaxDiffIndex(pts mat.Dense) (int, int){
	r, c := pts.Dims()

	maxIndex := 0
	maxValue := 0.0
	minIndex := 0
	minValue := 0.0

	for i := 0; i < r; i++ {
		row := pts.RowView(i)
		diff := row.AtVec(c - 1) //Do check c is not Zero or 1
		for j := c - 2; j >= 0; j-- {
			diff -= row.AtVec(j)
		}

		if (i == 0 ) {
			maxValue = diff
			minValue = diff
		}

		//Find max value and index
		if (diff > maxValue) {
			maxValue = diff
			maxIndex = i
		}
		//Find min value and index
		if (diff < minValue) {
			minValue = diff
			minIndex = i
		}
	}
	return minIndex, maxIndex
}

func convertDenseToImagePoint(pts *mat.Dense) []image.Point {
	var sd []image.Point

	r, c := pts.Dims()
	if (c !=2 ) {
		return sd
	}
	for i := 0; i < r; i++ {
		row := pts.RowView(i)
		sd = append(sd, image.Point{
			X: int(row.AtVec(0)),
			Y: int(row.AtVec(1)),
		})
	}
	return sd
}
