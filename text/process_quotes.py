quotes = []
cur = []

for line in open("quotes.txt", encoding='latin1'):
    if line.startswith('    '):
        cur.append(line.strip())
    elif cur:
        quotes.append(' '.join(cur))
        cur = []
if cur:
    quotes.append(' '.join(cur))

quotes = filter(lambda x: x.isascii(), quotes)
quotes = filter(lambda x: 40 <= len(x) <= 80, quotes)
for quote in quotes:
    print(quote)
