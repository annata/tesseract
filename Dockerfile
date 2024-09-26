FROM golang:1.22-alpine as builder
WORKDIR /ocr
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -ldflags="-w -s"

FROM buildpack-deps:curl
RUN apt update && apt install -y ca-certificates curl gnupg lsb-release && echo "deb https://notesalexp.org/tesseract-ocr5/$(lsb_release -cs)/ $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/notesalexp.list > /dev/null && wget -O - https://notesalexp.org/debian/alexp_key.asc | apt-key add - && apt-get update && apt-get install tesseract-ocr tesseract-ocr-chi-sim tesseract-ocr-chi-sim-best tesseract-ocr-eng-best -y && apt-get remove --purge --auto-remove -y && rm -rf /var/lib/apt/lists/*
ENV TESSDATA_PREFIX /usr/share/tesseract-ocr/5/tessdata/
COPY --from=builder /ocr/ocr /
CMD ["/ocr"]