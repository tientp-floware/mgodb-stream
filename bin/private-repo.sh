## Enable private repo
go env -w GOPRIVATE=*.ghn.vn
## Set git config
## Example
## USER: tientp, PERSONAL_TOKEN: secret string
git config --global \
  url."https://${USER}:${PERSONAL_TOKEN}@g.ghn.vn/".insteadOf \
  "https://g.ghn.vn/"