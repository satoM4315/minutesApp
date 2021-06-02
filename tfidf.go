package main

//引用先：https://github.com/ramenjuniti/jtfidf/blob/master/jtfidf.go#L77
import (
	"github.com/ikawaha/kagome/tokenizer"
	"math"
)

func splitTerm(d string) []string {
	t := tokenizer.New()
	tokens := t.Tokenize(d)
	tokens = tokens[1 : len(tokens)-1]
	terms := make([]string, len(tokens))

	for i, token := range tokens {
		terms[i] = token.Surface
	}

	return terms
}

// AllTf returns all TF values in d
func AllTf(d string) map[string]float64 {
	terms := splitTerm(d)
	n := len(terms)
	tfs := map[string]float64{}

	for _, term := range terms {
		if _, ok := tfs[term]; ok {
			tfs[term]++
		} else {
			tfs[term] = 1
		}
	}

	for term := range tfs {
		tfs[term] /= float64(n)
	}

	return tfs
}

// Tf returns t's TF value in d
func Tf(t, d string) float64 {
	terms := splitTerm(d)
	n := len(terms)
	var count int

	if n == 0 {
		return 0
	}

	for _, term := range terms {
		if t == term {
			count++
		}
	}

	return float64(count) / float64(n)
}

// AllIdf returns all IDF values in ds
func AllIdf(ds []string) map[string]float64 {
	n := len(ds)
	terms := []string{}
	termsList := make([][]string, n)

	for _, d := range ds {
		terms = append(terms, splitTerm(d)...)
	}

	for i, d := range ds {
		termsList[i] = splitTerm(d)
	}

	idfs := map[string]float64{}

	for _, term := range terms {
		var df int
		for i := 0; i < len(termsList); i++ {
			for j := 0; j < len(termsList[i]); j++ {
				if termsList[i][j] == term {
					df++
					break
				}
			}
		}
		if _, ok := idfs[term]; !ok {
			idfs[term] = math.Log(float64(n) / float64(df))
		}
	}

	return idfs
}

// Idf retuns t's IDF value in ds
func Idf(t string, ds []string) float64 {
	n := len(ds)
	termsList := make([][]string, n)
	var df int

	for i, d := range ds {
		termsList[i] = splitTerm(d)
	}

	for i := 0; i < len(termsList); i++ {
		for j := 0; j < len(termsList[i]); j++ {
			if t == termsList[i][j] {
				df++
				break
			}
		}
	}

	if df == 0 {
		return 0
	}

	return math.Log(float64(n) / float64(df))
}

// AllTfidf retuns all TF-IDF values in ds
func allTfIdf(ds []string) []map[string]float64 {
	idfs := AllIdf(ds)
	tfidfs := []map[string]float64{}

	for _, d := range ds {
		tfidf := map[string]float64{}
		for term, tf := range AllTf(d) {
			tfidf[term] = tf * (idfs[term] + 1)
		}
		tfidfs = append(tfidfs, tfidf)
	}

	return tfidfs
}

// Tfidf returns t's TF-IDF value in ds
func Tfidf(t, d string, ds []string) float64 {
	return Tf(t, d) * (Idf(t, ds) + 1)
}

func allTfIdfVec(ds []string) [][]float64 {
	tfidfs := allTfIdf(ds)
	idfs := AllIdf(ds)
	vocab := make([]string, len(idfs))
	index := 0
	for term, _ := range idfs {
		vocab[index] = term
		index++
	}
	tfidfVec := make([][]float64, len(ds))
	for i, tfidf := range tfidfs {
		vec := make([]float64, len(vocab))
		for _, term := range vocab {
			if _, ok := tfidf[term]; ok {
				vec = append(vec, tfidf[term])
			} else {
				vec = append(vec, 0)
			}
		}
		tfidfVec[i] = vec
	}
	return tfidfVec
}
