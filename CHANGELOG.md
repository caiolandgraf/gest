# Changelog

## [v2.0.0-beta]

### Refactor
- Mudança do pacote `gest` para `package main`, tornando-o uma ferramenta de linha de comando.

### New Features
- Implementação de parser para o output JSON de `go test`.
- Integração com `fsnotify` para modo de execução interativo/watch.

### Breaking Changes
- Remoção da API interna de asserções em favor do parsing de resultados de testes existentes.

### Update
- Ajuste nos exemplos para refletir a nova estrutura de importação.
