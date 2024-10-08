# Kensaku

Keansaku is the japanese word for "search".

## Quickstart

```bash
rm -f db/kensaku.db
sqlite3 db/kensaku.db < db/schema.sql
go run . serve
```

Then point your browser to `https://localhost:8080`

## Indexing text documents

```bash
go run . index test_files/example1.txt
go run . index test_files/example2.txt
go run . index test_files/example3.txt
go run . index test_files/example4.txt
go run . index test_files/news.html
```

## Indexing URLs and HTML files

TBD

## Indexing Images

TBD

## Indexing Office files

TBD

## Importing Objects from CSV files

TBD

## Indexing PDF files

TBD

## Indexing password protected URLs

TBD

## Indexing password protected files

TBD

## Indexing compressed files?

TBD
