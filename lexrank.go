package main

import (
	"gonum.org/v1/gonum/mat"
	"math"
	"sort"
)

func getImportantSentence(ds []string) []int {
	tfidfVec := allTfIdfVec(ds)
	numDoc := len(ds)
	forMatrix := make([]float64, numDoc*numDoc)
	for i, vec1 := range tfidfVec {
		for j, vec2 := range tfidfVec {
			forMatrix[i*numDoc+j] = cosSim(vec1, vec2)
		}
	}
	simMatrix := mat.NewDense(numDoc, numDoc, forMatrix)
	thr := 0.1
	adjacencyMatrix := mat.NewDense(numDoc, numDoc, nil)
	for i := 0; i < numDoc; i++ {
		for j := 0; j < numDoc; j++ {
			if simMatrix.At(i, j) > thr {
				adjacencyMatrix.Set(i, j, 1)
			}
		}
	}
	stochasticMatrix := adjacencyMatrix
	for i := 0; i < numDoc; i++ {
		degree := mat.Sum(stochasticMatrix.Slice(i, i+1, 0, numDoc))
		if degree > 0 {
			for j := 0; j < numDoc; j++ {
				stochasticMatrix.Set(i, j, stochasticMatrix.At(i, j)/degree)
			}
		}
	}
	forP := make([]float64, numDoc)
	for i := 0; i < numDoc; i++ {
		forP[i] = 1.0 / float64(numDoc)
	}
	p := mat.NewDense(numDoc, 1, forP)
	delta := 1.0
	for delta > 10e-6 {
		_p := mat.NewDense(numDoc, 1, nil)
		_p.Product(stochasticMatrix.T(), p)
		sub := mat.NewDense(numDoc, 1, nil)
		sub.Sub(_p, p)
		delta = mat.Norm(sub, 2)
		p = _p
	}

	result := make([]float64, numDoc)
	for i := 0; i < numDoc; i++ {
		result[i] = p.At(i, 0)
	}
	type Lexrank struct {
		score float64
		index int
	}
	lexRank := make([]Lexrank, numDoc)
	for i, value := range result {
		lexRank[i] = Lexrank{value, i}
	}
	sort.Slice(lexRank, func(i, j int) bool { return lexRank[i].score > lexRank[j].score })
	ranking := make([]int, numDoc)
	for i, rank := range lexRank {
		ranking[i] = rank.index
	}
	return ranking
}

func cosSim(vec1 []float64, vec2 []float64) (ret float64) {
	ret = dot(vec1, vec2)
	ret = ret / (length(vec1) * length(vec2))
	return
}

func dot(vec1 []float64, vec2 []float64) (ret float64) {
	for i := 0; i < len(vec1); i++ {
		ret += vec1[i] * vec2[i]
	}
	return
}

func length(vec []float64) (ret float64) {
	ret = math.Sqrt(dot(vec, vec))
	return
}
