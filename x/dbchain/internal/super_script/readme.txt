Since this package is a bit complicated, I added this readme file to help re-reading
the code in the future.

The purpose of this package is to implement a parser for super script.

0. The super script's syntax is depicted in file bnf.go
1. It scans script to tokens
2. It construct systax tree out of the tokens 
3. It parser.err represent whether the script is valid
4. If the script is valid, the parser.systaxTree is the result

5. In the ./eval pacakge, we use the systax tree to evaluate the script
   against the insertion operation
