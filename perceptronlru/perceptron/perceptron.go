/*
Copyright 2013 Alexandre Passos

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package perceptron

import "math"

type Perceptron struct {
	size    int32
	weights []float64
	adagrad []float64
}

func New(size int32) *Perceptron {
	return &Perceptron{
		size:    size,
		weights: make([]float64, size, size),
		adagrad: make([]float64, size, size),
	}
}

func sign(a int32) float64 {
	if a > 0 {
		return 1.0
	} else {
		return -1.0
	}
}

func abs(a int32) int32 {
	if a < 0 {
		a = -a
	}
	return a
}

func (model *Perceptron) Score(features []int32) float64 {
	score := 0.0
	for i := 0; i < len(features); i++ {
		n := features[i]
		score += sign(n) * model.weights[abs(n)%model.size]
	}
	return score
}

func (model *Perceptron) Update(features []int32, target float64) float64 {
	score := model.Score(features)
	var direction float64
	if target > score {
		direction = -1.0
	} else {
		direction = 1.0
	}
	for i := 0; i < len(features); i++ {
		n := features[i]
		grad := direction * sign(n)
		idx := abs(n) % model.size
		model.adagrad[idx] += grad * grad
		lrate := 1.0 / math.Sqrt(model.adagrad[idx])
		score += direction * lrate
		model.weights[idx] += grad * lrate
	}
	return score
}
