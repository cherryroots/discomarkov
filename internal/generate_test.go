package internal

import "testing"

func TestGetNgram(t *testing.T) {
	messages := []string{
		"Hello world",
		"Hello universe",
		"Hello testing world from the good people of the world",
	}
	ngram2, err := GetNgram(messages, 2)
	if err != nil {
		t.Fatal(err)
	}
	if ngram2.N != 2 {
		t.Errorf("Expected ngram2 N to be 2, got %d", ngram2.N)
	}
	if len(ngram2.Grams) != 11 {
		t.Errorf("Expected ngram2 Grams to have 11 elements, got %d", len(ngram2.Grams))
	}

	ngram3, err := GetNgram(messages, 3)
	if err != nil {
		t.Fatal(err)
	}
	if ngram3.N != 3 {
		t.Errorf("Expected ngram3 N to be 3, got %d", ngram3.N)
	}
	if len(ngram3.Grams) != 8 {
		t.Errorf("Expected ngram3 Grams to have 8 elements, got %d", len(ngram3.Grams))
	}
}

func TestGetWordFrequency(t *testing.T) {
	messages := []string{
		"Hello world",
		"Hello universe",
	}
	wordFrequency, err := GetWordFrequency(messages)
	if err != nil {
		t.Fatal(err)
	}
	if len(wordFrequency) != 3 {
		t.Errorf("Expected wordFrequency to have 3 elements, got %d", len(wordFrequency))
	}
	if wordFrequency["Hello"] != 2 {
		t.Errorf("Expected wordFrequency[\"Hello\"] to be 2, got %d", wordFrequency["Hello"])
	}
}
