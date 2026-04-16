# Issue (plantilla): quitar gorilla/mux

**GitHub CLI no estaba disponible en el entorno donde se generó el PR.** Crea el issue manualmente en  
https://github.com/gaston-garcia-cegid/GonsGarage/issues/new  
y pega lo siguiente (luego enlaza el issue en el PR si aplica).

## Título sugerido

`chore: eliminar gorilla/mux (user_handler migrado a Gin)`

## Cuerpo sugerido

```markdown
## Contexto
`backend/internal/adapters/http/handlers/user_handler.go` usaba **gorilla/mux** (`mux.Vars`) en un stub de usuarios.

## Objetivo
Eliminar la dependencia `github.com/gorilla/mux` del módulo Go y alinear el stub con **Gin** (`c.Param`), como el resto de handlers.

## Criterios de cierre
- [x] Sin import de `gorilla/mux` en el backend
- [x] `go mod tidy` sin `gorilla/mux` en `go.mod`

## Nota
El cambio puede ir en el mismo PR que migra el handler; este issue sirve para trazabilidad y búsqueda.
```
