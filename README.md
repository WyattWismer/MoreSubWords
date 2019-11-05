# More Subwords

**More Subwords** is a game where you are given a phrase and must find *subwords* in that phrase.

You can play the game [here](http://moresubwords.herokuapp.com/)

## Instructions
We define a *subword* as any word that:
1. Is a substring of the phrase
2. Is an english word

### Example:
Consider the phrase:
> A sailing ship goes west at night  

Some possible substrings include:
+ 'sail'
+ 'ship goes west'
+ 'h'
+ 'assgwan'
+ 'wet'
+ 'sat'
+ 'tight'
+ 'ships'
+ ...

In general, a substring is any string that can be made by deleting characters in the original string.

Out of these substrings the english words were:
+ 'sail'
+ 'wet'
+ 'sat'
+ 'tight'
+ 'ships'

To make the game more challenging  you only score points for subwords that:
1. Are not a prefix of an existing word
2. Do not have an existing word as their prefix

So in the example above we cannot get points for:
+ 'sail' since it is a prefix of 'sailing'
+ 'ships' since 'ship' is a prefix of ship

## Sources for Words & Phrases

Words and phrases used by this project are in the public domain.

Phrases taken from:
http://www.gutenberg.org/ebooks/27889

Words taken from:
https://github.com/dwyl/english-words
