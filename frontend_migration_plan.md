# 🚀 Plan de Migración Frontend - GonsGarage

## 📋 **RESUMEN EJECUTIVO**
- **Tiempo Total Estimado**: 8-12 horas de desarrollo
- **Duración Calendario**: 2-3 días de trabajo
- **Tareas Críticas**: 6 tareas principales
- **Dependencias**: 4 cadenas de dependencias identificadas
- **Testing**: Cada tarea incluye pruebas inmediatas

---

## 🎯 **FASE 1: FUNDACIÓN (Prioridad CRÍTICA)**
> **Objetivo**: Establecer base sólida para migración
> **Tiempo Total**: 3-4 horas

### **TAREA 1: Instalar Dependencies Core** ⭐⭐⭐⭐⭐
- **Tiempo**: 30-45 minutos
- **Prioridad**: CRÍTICA (Bloquea todo lo demás)
- **Dependencias**: Ninguna (PRIMERA TAREA)
- **Archivos a Modificar**: `package.json`

#### Subtareas:
```bash
# 1. Instalar Zustand ecosystem (15 min)
pnpm add zustand immer

# 2. Instalar testing dependencies (15 min)
pnpm add -D jest @testing-library/react @testing-library/jest-dom @testing-library/user-event

# 3. Configurar Jest (15 min)
# Crear jest.config.js y setup tests
```

#### **✅ Prueba Inmediata**:
```bash
# Verificar instalación
pnpm list zustand
pnpm list jest

# Verificar que no hay errores de dependencias
pnpm install
```

---

### **TAREA 2: Crear Estructura Base de Types** ⭐⭐⭐⭐⭐
- **Tiempo**: 45-60 minutos
- **Prioridad**: CRÍTICA
- **Dependencias**: Ninguna (Paralela con Tarea 1)
- **Archivos a Crear**: 
  - `src/types/index.ts`
  - `src/types/auth.ts`
  - `src/types/user.ts`
  - `src/types/employee.ts`
  - `src/types/api.ts`

#### **✅ Prueba Inmediata**:
```bash
# Verificar compilación TypeScript
pnpm run build

# O verificar type checking
npx tsc --noEmit
```

---

### **TAREA 3: Configurar Testing Framework** ⭐⭐⭐⭐
- **Tiempo**: 60-75 minutos
- **Prioridad**: ALTA
- **Dependencias**: TAREA 1 (necesita jest instalado)
- **Archivos a Crear**:
  - `jest.config.js`
  - `__tests__/setup.ts`
  - `__tests__/example.test.ts`

#### **✅ Prueba Inmediata**:
```bash
# Ejecutar test de ejemplo
pnpm test

# Verificar coverage funciona
pnpm test -- --coverage
```

---

## 🏗️ **FASE 2: MIGRACIÓN CORE (Prioridad ALTA)**
> **Objetivo**: Migrar estado y API client
> **Tiempo Total**: 3-4 horas

### **TAREA 4: Crear AuthStore (Zustand)** ⭐⭐⭐⭐⭐
- **Tiempo**: 90-120 minutos
- **Prioridad**: CRÍTICA
- **Dependencias**: TAREA 1 + TAREA 2 (necesita types y zustand)
- **Archivos a Crear**:
  - `src/stores/auth.store.ts`
  - `src/stores/index.ts`
- **Archivos a Analizar**: `src/contexts/AuthContext.tsx`

#### Subtareas:
1. **Analizar AuthContext actual** (20 min)
2. **Crear AuthStore con Zustand** (40 min)
3. **Mantener misma API** (30 min)
4. **Crear tests básicos** (30 min)

#### **✅ Prueba Inmediata**:
```bash
# Test del store
pnpm test -- auth.store.test.ts

# Importar store en componente de prueba
# Verificar que compila sin errores
pnpm run build
```

---

### **TAREA 5: Migrar API Client Structure** ⭐⭐⭐⭐
- **Tiempo**: 90-120 minutos  
- **Prioridad**: ALTA
- **Dependencias**: TAREA 2 (necesita types/api.ts)
- **Archivos Actual**: `src/lib/api.ts`
- **Archivos Nuevos**:
  - `src/lib/api/index.ts`
  - `src/lib/api/auth.api.ts`
  - `src/lib/api/user.api.ts`
  - `src/lib/api/car.api.ts`

#### **✅ Prueba Inmediata**:
```bash
# Test de cada API client
pnpm test -- auth.api.test.ts
pnpm test -- car.api.test.ts

# Verificar imports funcionan
pnpm run build
```

---

### **TAREA 6: Migrar Componentes de AuthContext → AuthStore** ⭐⭐⭐⭐
- **Tiempo**: 60-90 minutos
- **Prioridad**: ALTA  
- **Dependencias**: TAREA 4 (necesita AuthStore creado)
- **Archivos a Modificar**: Todos los que usan AuthContext

#### **✅ Prueba Inmediata**:
```bash
# Test de componentes migrados
pnpm test -- components/

# Verificar app funciona
pnpm run dev
# Probar login/logout en navegador
```

---

## 🎨 **FASE 3: STORES ADICIONALES (Prioridad MEDIA)**
> **Objetivo**: Completar ecosystem Zustand  
> **Tiempo Total**: 2-3 horas

### **TAREA 7: Crear UserStore y CarStore** ⭐⭐⭐
- **Tiempo**: 90-120 minutos
- **Prioridad**: MEDIA
- **Dependencias**: TAREA 4 + TAREA 5 (patrón establecido)
- **Archivos a Crear**:
  - `src/stores/user.store.ts`
  - `src/stores/car.store.ts`

#### **✅ Prueba Inmediata**:
```bash
# Test de cada store
pnpm test -- user.store.test.ts
pnpm test -- car.store.test.ts

# Test de integración con API
pnpm test -- integration/
```

---

### **TAREA 8: Migrar Hooks a Zustand** ⭐⭐⭐
- **Tiempo**: 45-60 minutos
- **Prioridad**: MEDIA
- **Dependencias**: TAREA 7 (necesita stores completos)
- **Archivos a Modificar**:
  - `src/hooks/useCars.ts`
  - Otros hooks custom

#### **✅ Prueba Inmediata**:
```bash
# Test de hooks migrados
pnpm test -- hooks/

# Verificar componentes que usan hooks
pnpm run dev
```

---

## 🧪 **FASE 4: TESTING COMPLETO (Prioridad BAJA)**
> **Objetivo**: Coverage completo y CI/CD
> **Tiempo Total**: 1-2 horas

### **TAREA 9: Completar Test Coverage** ⭐⭐
- **Tiempo**: 60-90 minutos
- **Prioridad**: BAJA
- **Dependencias**: Todas las anteriores
- **Objetivo**: >80% coverage

#### **✅ Prueba Inmediata**:
```bash
# Coverage completo
pnpm test -- --coverage

# Verificar threshold
pnpm test -- --coverage --watchAll=false
```

---

## 🔄 **CADENAS DE DEPENDENCIAS IDENTIFICADAS**

### **Cadena A: Base Foundation**
```
TAREA 1 (Dependencies) 
    ↓
TAREA 3 (Testing Setup)
    ↓  
TAREA 4 (AuthStore)
```

### **Cadena B: Types & API**
```
TAREA 2 (Types)
    ↓
TAREA 5 (API Client)
    ↓
TAREA 7 (Additional Stores)
```

### **Cadena C: Components Migration**
```
TAREA 4 (AuthStore)
    ↓
TAREA 6 (Component Migration)
    ↓
TAREA 8 (Hooks Migration)
```

### **Cadena D: Testing Complete**
```
TAREA 7 + TAREA 8
    ↓
TAREA 9 (Complete Coverage)
```

---

## 📅 **CRONOGRAMA SUGERIDO**

### **DÍA 1 (4-5 horas)**
- **Mañana**: TAREA 1 + TAREA 2 + TAREA 3 (Base foundation)
- **Tarde**: TAREA 4 (AuthStore) + inicio TAREA 5

### **DÍA 2 (3-4 horas)**  
- **Mañana**: Completar TAREA 5 + TAREA 6 (Migration)
- **Tarde**: TAREA 7 (Additional Stores)

### **DÍA 3 (1-2 horas)**
- **Mañana**: TAREA 8 + TAREA 9 (Final testing)

---

## 🚦 **CRITERIOS DE ÉXITO POR FASE**

### **FASE 1 COMPLETA CUANDO**:
- ✅ `pnpm install` sin errores
- ✅ `pnpm run build` exitoso
- ✅ `pnpm test` ejecuta sin errores
- ✅ Types exportados correctamente

### **FASE 2 COMPLETA CUANDO**:
- ✅ AuthStore reemplaza AuthContext completamente
- ✅ API clients funcionan en tests
- ✅ Login/logout funciona en navegador
- ✅ No quedan referencias a AuthContext

### **FASE 3 COMPLETA CUANDO**:
- ✅ Todos los stores implementados  
- ✅ Hooks migrados a Zustand
- ✅ CRUD operations funcionan

### **FASE 4 COMPLETA CUANDO**:
- ✅ Coverage > 80%
- ✅ CI/CD pasa todos los tests
- ✅ Documentación actualizada

---

## ❓ **¿POR DÓNDE EMPEZAMOS?**

**Recomendación**: Comenzar con **TAREA 1** inmediatamente porque:

1. **Bloquea todo lo demás** - Sin Zustand no podemos crear stores
2. **Se puede probar inmediatamente** - Instalación verificable
3. **No tiene dependencias** - Podemos empezar ahora mismo
4. **Tiempo corto** - 30-45 minutos máximo

¿Te parece bien empezar con la instalación de dependencies?