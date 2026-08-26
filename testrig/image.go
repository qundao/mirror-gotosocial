// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package testrig

import (
	"bytes"
	"image"
	"math"
	"testing"
)

// maxMeanAbsDiffPerChannel is the maximum allowed mean difference per channel (0-255)
// when comparing two JPEGs for similarity.
const maxMeanAbsDiffPerChannel = 5

// AssertSimilarImages loads two images and asserts they have the same dimensions and
// similar enough content to be considered the same.
func AssertSimilarImages(t *testing.T, testImageBytes, refImageBytes []byte) {
	t.Helper()

	imgTest, _, err := image.Decode(bytes.NewReader(testImageBytes))
	if err != nil {
		t.Fatalf("test image is not valid: %v", err)
	}
	imgRef, _, err := image.Decode(bytes.NewReader(refImageBytes))
	if err != nil {
		t.Fatalf("reference image is not valid: %v", err)
	}

	bTest := imgTest.Bounds()
	bRef := imgRef.Bounds()
	if bTest.Dx() != bRef.Dx() || bTest.Dy() != bRef.Dy() {
		t.Fatalf("dimension mismatch: test %dx%d, reference %dx%d", bTest.Dx(), bTest.Dy(), bRef.Dx(), bRef.Dy())
	}

	// Calculate the mean difference in each channel between the test and reference images
	var sumR, sumG, sumB, sumA float64
	w, h := bTest.Dx(), bTest.Dy()
	n := float64(w * h)
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			r1, g1, b1, a1 := imgTest.At(bTest.Min.X+i, bTest.Min.Y+j).RGBA()
			r2, g2, b2, a2 := imgRef.At(bRef.Min.X+i, bRef.Min.Y+j).RGBA()
			// RGBA() returns an unsigned 16-bit integer
			// convert to signed int (0-255) and get the absolute difference
			sumR += math.Abs(float64(int(r1>>8) - int(r2>>8)))
			sumG += math.Abs(float64(int(g1>>8) - int(g2>>8)))
			sumB += math.Abs(float64(int(b1>>8) - int(b2>>8)))
			sumA += math.Abs(float64(int(a1>>8) - int(a2>>8)))
		}
	}
	meanR := sumR / n
	meanG := sumG / n
	meanB := sumB / n
	meanA := sumA / n
	if meanR > maxMeanAbsDiffPerChannel || meanG > maxMeanAbsDiffPerChannel || meanB > maxMeanAbsDiffPerChannel || meanA > maxMeanAbsDiffPerChannel {
		t.Errorf(
			"images differ too much: mean diff per channel = R=%f G=%f B=%f A=%f (max %d)",
			meanR, meanG, meanB, meanA, maxMeanAbsDiffPerChannel,
		)
	}
}
