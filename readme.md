# DNA Sequencing tools

Tools to help sequence and visualize binary and non-binary files.

The nucleotide in DNA consists of a sugar (deoxyribose), one of four bases (cytosine (C), thymine (T), adenine (A), guanine (G)), and a phosphate. Cytosine and thymine are pyrimidine bases, while adenine and guanine are purine bases. The sugar and the base together are called a nucleoside.

**Table of contents**

- [Included tools](#included-tools)
  - [dna-encoder](#dna-encoder)
  - [fasta-to-image](#fasta-to-image)
- [Example with normal text file](#example-with-normal-text-file)

## Included tools

- dna-encoder
- fasta-to-image

### dna-encoder

Takes a file (binary or non-binary) and encodes it to DNA sequence.

Binary representation of every byte is translated into DNA by encoding the following.

```
00    A    Adenine     color.RGBA{0, 0, 255, 255}
01    G    Guanine     color.RGBA{0, 100, 0, 255}
10    C    Cytosine    color.RGBA{255, 0, 0, 255}
11    T    Thymine     color.RGBA{255, 255, 0, 255}
```

[![asciicast](https://asciinema.org/a/EvfFa4n7Cr9DVzbw4HNaY5323.png | width=100)](https://asciinema.org/a/EvfFa4n7Cr9DVzbw4HNaY5323)

### fasta-to-image

Takes FASTA file which is outputed from dna-encoder and creates PNG image.

[![asciicast](https://asciinema.org/a/VJmaBdoYp5sgZelqESi96MU4S.svg)](https://asciinema.org/a/VJmaBdoYp5sgZelqESi96MU4S)

## Example with normal text file

**Original**

```
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis et consectetur turpis. Integer quis pharetra turpis. Donec dui mauris, dignissim eu elementum nec, euismod id orci.
```

**FASTA file**

```fasta
>SEQ1
GAAGCGACGCGGGCGACAAGCCGGAAGAGGGGCGACAAGCGAGCGCAGCGACACAAGAGC
CGGGAACAAGCAGGCGGCGGGGAACAACAAGCAGCGCCGAGCGGGCAGGAGCGGGGAGGG
GACACAAGCAGGCGAGCCGGAAGCCGGAGCAGCCGGCCGCGACAAGCGGGCAGCCGGGAA
CCACAAGAGAGGGGCCGGAACAAGCGGGGAACAAGCAGCGCCGAGCGGGCAGGAGCGGGG
AGGGGACACAAGGAGGGGACGAAGCCGGAACCACAAGACGGCCGGAGCGGGCGGCGGGAC
ACAAGAGGGGGCCGGAACAAGAAGCCAGCAGGACGCGGGGAGACGCAGACAAGGAGGGGA
CGAAGCCGGAACCACAAGAGAGCGCCGCGGGCAACAAGCGAGGGGCCGACAAGCGGCAGG
GGGACGCCGGAACAACAAGCGAGCCGGCGGCCGCCGGAGAGCCGGCGACAAGCGGGGGAC
AAGCGGGCAGCGGGCGGCGGGCCGGAGGGGCGACAAGCCGCGGGCAACAACAAGCGGGGG
GCCGGAGCGGCGCGAACAAGCCGGCGAACAAGCGACGCAGCCGACC
```

**Encoded into image**

![dna](https://user-images.githubusercontent.com/296714/50626024-22e63280-0f2c-11e9-8d86-7f75d35b1804.png)
