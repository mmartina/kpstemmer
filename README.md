KP-Stemmer package for Go
=========================
Implements the Kraaij-Pohlmann stemmer as described at
   <http://snowball.tartarus.org/algorithms/kraaij_pohlmann/stemmer.html>

Installation
-------------
    go get github.com/mmartina/kpstemmer

Example
-------
    import "github.com/mmartina/kpstemmer"

    kpstemmer.Stem("lichamelijke") // => lichamelijk
    kpstemmer.Stem("opglimpende")  // => opglimp

Tests
-----
Included `test_diffs.txt` as listed at
    <http://snowball.tartarus.org/algorithms/kraaij_pohlmann/diffs.txt>

Used only when running tests with `go test`.

License
-------
MIT License (see LICENSE file).