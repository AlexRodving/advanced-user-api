# 🌐 Documentation Website Guide

## Обзор

Сайт документации создан с помощью **Docsify** - генератора документации без сборки, работающего прямо с Markdown файлами.

**URL сайта:** https://alexrodving.github.io/advanced-user-api/

---

## 🛠️ Технологии

### Docsify

**Что это:** Генератор документации, работающий в браузере

**Официальный сайт:** https://docsify.js.org/

**Преимущества:**
- ✅ Не требует сборки (build)
- ✅ Работает прямо с Markdown
- ✅ Динамическая загрузка страниц
- ✅ Множество плагинов
- ✅ Настраиваемый дизайн
- ✅ SEO-friendly

**Альтернативы:**
- MkDocs (Python, требует сборку)
- Docusaurus (React, требует сборку)
- VitePress (Vue, требует сборку)

---

## 📁 Структура сайта

```
docs/                         # GitHub Pages папка
├── index.html               # Главная страница + конфигурация Docsify
├── README.md                # Домашняя страница сайта
├── _sidebar.md              # Конфигурация бокового меню
├── _navbar.md               # Конфигурация верхнего меню
├── .nojekyll                # Отключение Jekyll обработки
└── docs/                    # Вся документация проекта
    ├── API.md
    ├── ARCHITECTURE.md
    ├── TESTING.md
    ├── GIT_WORKFLOW.md
    ├── QUICKSTART.md
    ├── DEPLOY.md
    ├── PROJECT_SUMMARY.md
    ├── WEBSITE.md           # Этот файл
    └── libraries/
        ├── README.md
        ├── GIN.md
        ├── GORM.md
        ├── JWT.md
        └── BCRYPT_VIPER.md
```

---

## 📄 Файлы сайта

### 1. index.html

**Назначение:** Главная страница и конфигурация Docsify

**Структура:**
```html
<!DOCTYPE html>
<html>
<head>
  <!-- Meta tags -->
  <meta charset="UTF-8">
  <title>Advanced User API - Documentation</title>
  
  <!-- Docsify theme -->
  <link rel="stylesheet" href="//cdn.jsdelivr.net/npm/docsify@4/lib/themes/vue.css">
  
  <!-- Custom CSS -->
  <style>
    /* Кастомные стили */
  </style>
</head>
<body>
  <div id="app">Загрузка...</div>
  
  <script>
    window.$docsify = {
      /* Конфигурация Docsify */
    }
  </script>
  
  <!-- Scripts -->
  <script src="//cdn.jsdelivr.net/npm/docsify@4"></script>
  <!-- Плагины -->
</body>
</html>
```

**Что включает:**

#### CSS стили:
```css
/* Gradient Hero */
.hero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 60px 20px;
  border-radius: 0 0 20px 20px;
}

/* Tech Stack Icons */
.tech-stack {
  display: flex;
  gap: 30px;
}

.tech-item:hover {
  transform: translateY(-5px);  /* Эффект поднятия при hover */
}

/* Feature Cards */
.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.feature-card {
  background: #f8f9fa;
  border-left: 4px solid var(--theme-color);
  transition: all 0.3s;
}

.feature-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
}

/* Stats Counters */
.stats {
  display: flex;
  justify-content: space-around;
}

.stat-number {
  font-size: 2.5em;
  font-weight: bold;
  color: var(--theme-color);
}

/* Alert Boxes */
.alert-info {
  background: #e3f2fd;
  border-color: #2196f3;
}
```

**Файл:** `docs/index.html`

---

### 2. Docsify Configuration

```javascript
window.$docsify = {
  // Название проекта
  name: 'Advanced User API',
  
  // Ссылка на GitHub
  repo: 'https://github.com/AlexRodving/advanced-user-api',
  
  // Загрузка sidebar и navbar
  loadSidebar: true,
  loadNavbar: true,
  
  // Уровень вложенности заголовков в sidebar
  subMaxLevel: 3,
  
  // Автоскролл вверх при переходе
  auto2top: true,
  
  // Настройки поиска
  search: {
    placeholder: 'Поиск...',
    noData: 'Ничего не найдено',
    depth: 6                    // Глубина поиска
  },
  
  // Копирование кода
  copyCode: {
    buttonText: 'Копировать',
    successText: 'Скопировано!'
  },
  
  // Pagination (переход между страницами)
  pagination: {
    previousText: 'Предыдущая',
    nextText: 'Следующая',
    crossChapter: true
  }
}
```

---

### 3. README.md (главная страница)

**Назначение:** Содержимое главной страницы сайта

**Особенности:**

#### HTML элементы в Markdown:
```html
<div class="hero">
  <h1>🚀 Advanced User API</h1>
  <p>Production-ready REST API с JWT аутентификацией</p>
</div>
```

#### Иконки технологий (DevIcons CDN):
```html
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg">
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postgresql/postgresql-original.svg">
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg">
```

**Источник иконок:** https://devicon.dev/

#### Feature Cards:
```html
<div class="features-grid">
  <div class="feature-card">
    <h3>🔐 JWT Authentication</h3>
    <p>Описание...</p>
  </div>
</div>
```

#### Statistics:
```html
<div class="stats">
  <div class="stat-item">
    <div class="stat-number">2,300+</div>
    <div class="stat-label">Строк кода</div>
  </div>
</div>
```

**Файл:** `docs/README.md`

---

### 4. _sidebar.md (боковое меню)

**Назначение:** Конфигурация навигации в левом sidebar

**Формат:**
```markdown
- **Раздел 1**
  - [Страница 1](docs/page1.md)
  - [Страница 2](docs/page2.md)

- **Раздел 2**
  - [Страница 3](docs/page3.md)
```

**Пример из проекта:**
```markdown
- **Начало**
  - [Главная](/)
  - [Quick Start](docs/QUICKSTART.md)

- **Документация**
  - [API](docs/API.md)
  - [Architecture](docs/ARCHITECTURE.md)

- **Библиотеки**
  - [Gin](docs/libraries/GIN.md)
  - [GORM](docs/libraries/GORM.md)
```

**Файл:** `docs/_sidebar.md`

---

### 5. _navbar.md (верхнее меню)

**Назначение:** Конфигурация навигации вверху страницы

**Формат:**
```markdown
- [🏠 Главная](/)
- [📡 API](docs/API.md)
- [GitHub](https://github.com/user/repo)
```

**Файл:** `docs/_navbar.md`

---

### 6. .nojekyll

**Назначение:** Отключает обработку Jekyll на GitHub Pages

**Почему нужен:**
- GitHub Pages по умолчанию использует Jekyll
- Jekyll игнорирует файлы начинающиеся с `_` (например `_sidebar.md`)
- `.nojekyll` отключает Jekyll и позволяет использовать все файлы

**Файл:** Пустой файл `docs/.nojekyll`

---

## 🎨 Дизайн компоненты

### Gradient Hero Section

```css
.hero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 60px 20px;
  text-align: center;
  border-radius: 0 0 20px 20px;
}
```

**Gradient источник:** https://uigradients.com/ (Piglet градиент)

**Можно изменить на:**
```css
/* Ocean Blue */
background: linear-gradient(135deg, #2E3192 0%, #1BFFFF 100%);

/* Sunset */
background: linear-gradient(135deg, #ff6e7f 0%, #bfe9ff 100%);

/* Green Beach */
background: linear-gradient(135deg, #02AAB0 0%, #00CDAC 100%);
```

---

### Tech Stack Icons

```html
<div class="tech-item">
  <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg" alt="Go">
  <div>Go 1.25</div>
</div>
```

**Доступные иконки:** https://devicon.dev/

**Другие полезные иконки:**
```html
<!-- Kubernetes -->
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/kubernetes/kubernetes-plain.svg">

<!-- GitHub Actions -->
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/github/github-original.svg">

<!-- Nginx -->
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nginx/nginx-original.svg">
```

---

### Feature Cards

```css
.feature-card {
  background: #f8f9fa;
  padding: 25px;
  border-radius: 12px;
  border-left: 4px solid var(--theme-color);
  transition: all 0.3s;
}

.feature-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
}
```

**Эффект:** Карточка поднимается и получает тень при наведении

---

### Alert Boxes

```css
.alert {
  padding: 15px 20px;
  margin: 20px 0;
  border-radius: 8px;
  border-left: 4px solid;
}

.alert-info {
  background: #e3f2fd;
  border-color: #2196f3;
  color: #0d47a1;
}
```

**Использование в Markdown:**
```html
<div class="alert alert-info">
💡 <strong>Совет:</strong> Начните с Quick Start Guide
</div>

<div class="alert alert-success">
✅ <strong>Успех:</strong> Сервер запущен
</div>

<div class="alert alert-warning">
⚠️ <strong>Внимание:</strong> Не используйте в production
</div>
```

---

## 🔌 Docsify Plugins

### Search Plugin

```html
<script src="//cdn.jsdelivr.net/npm/docsify/lib/plugins/search.min.js"></script>
```

**Конфигурация:**
```javascript
search: {
  placeholder: 'Поиск...',     // Плейсхолдер в поле
  noData: 'Ничего не найдено', // Если ничего не найдено
  depth: 6                      // Глубина поиска (уровни заголовков)
}
```

**Как работает:**
- Индексирует все Markdown файлы
- Поиск по заголовкам и тексту
- Показывает совпадения с подсветкой
- Переход на страницу по клику

---

### Copy Code Plugin

```html
<script src="//cdn.jsdelivr.net/npm/docsify-copy-code@2"></script>
```

**Конфигурация:**
```javascript
copyCode: {
  buttonText: 'Копировать',
  errorText: 'Ошибка',
  successText: 'Скопировано!'
}
```

**Как работает:**
- Добавляет кнопку "Копировать" к каждому блоку кода
- Копирует код в clipboard
- Показывает уведомление об успехе

---

### Pagination Plugin

```html
<script src="//cdn.jsdelivr.net/npm/docsify-pagination@2/dist/docsify-pagination.min.js"></script>
```

**Конфигурация:**
```javascript
pagination: {
  previousText: 'Предыдущая',
  nextText: 'Следующая',
  crossChapter: true           // Переход между разделами
}
```

**Как работает:**
- Добавляет кнопки "Предыдущая/Следующая" внизу страницы
- Автоматически определяет порядок из _sidebar.md

---

### Emoji Plugin

```html
<script src="//cdn.jsdelivr.net/npm/docsify/lib/plugins/emoji.min.js"></script>
```

**Как работает:**
- Преобразует :emoji: в 🎉
- Примеры: `:rocket:` → 🚀, `:fire:` → 🔥

---

### Syntax Highlighting (Prism)

```html
<script src="//cdn.jsdelivr.net/npm/prismjs@1/components/prism-bash.min.js"></script>
<script src="//cdn.jsdelivr.net/npm/prismjs@1/components/prism-go.min.js"></script>
<script src="//cdn.jsdelivr.net/npm/prismjs@1/components/prism-sql.min.js"></script>
<script src="//cdn.jsdelivr.net/npm/prismjs@1/components/prism-json.min.js"></script>
<script src="//cdn.jsdelivr.net/npm/prismjs@1/components/prism-yaml.min.js"></script>
<script src="//cdn.jsdelivr.net/npm/prismjs@1/components/prism-docker.min.js"></script>
```

**Поддерживаемые языки:**
- `bash` - Bash скрипты
- `go` - Go код
- `sql` - SQL запросы
- `json` - JSON данные
- `yaml` - YAML конфигурация
- `docker` - Dockerfile

**Использование в Markdown:**
````markdown
```go
func main() {
    fmt.Println("Hello")
}
```

```bash
docker compose up -d
```

```sql
SELECT * FROM users;
```
````

---

## 🎨 Кастомные CSS компоненты

### CSS Variables

```css
:root {
  --theme-color: #00ADD8;        /* Основной цвет (Go голубой) */
  --theme-color-dark: #007d9c;   /* Тёмный вариант */
  --sidebar-width: 280px;         /* Ширина sidebar */
}
```

**Как изменить тему:**
```css
:root {
  --theme-color: #764ba2;  /* Фиолетовый */
}
```

---

### Grid Layout

```css
.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}
```

**Как работает:**
- `repeat(auto-fit, ...)` - автоматическое количество колонок
- `minmax(280px, 1fr)` - минимум 280px, максимум равномерно
- `gap: 20px` - отступы между элементами

**На больших экранах:** 3 колонки  
**На планшетах:** 2 колонки  
**На мобильных:** 1 колонка

---

### Flexbox Layout

```css
.tech-stack {
  display: flex;
  justify-content: center;  /* Центрирование */
  flex-wrap: wrap;          /* Перенос на новую строку */
  gap: 30px;                /* Отступы */
}
```

---

### Hover Effects

```css
.feature-card:hover {
  transform: translateY(-3px);              /* Поднятие на 3px */
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);  /* Тень */
}
```

**Параметры transition:**
```css
transition: all 0.3s;  /* Все свойства, 0.3 секунды */
```

---

### Responsive Design

```css
@media (max-width: 768px) {
  .hero h1 {
    font-size: 2em;      /* Меньший размер на мобильных */
  }
  
  .tech-item img {
    width: 40px;         /* Меньшие иконки */
    height: 40px;
  }
}
```

**Breakpoints:**
- Desktop: > 768px
- Mobile: ≤ 768px

---

## 🖼️ Изображения

### Источники иконок

#### 1. DevIcons (используем)
```
https://cdn.jsdelivr.net/gh/devicons/devicon/icons/{name}/{name}-original.svg
```

**Примеры:**
- Go: `/go/go-original.svg`
- PostgreSQL: `/postgresql/postgresql-original.svg`
- Docker: `/docker/docker-original.svg`
- Redis: `/redis/redis-original.svg`

**Все иконки:** https://devicon.dev/

#### 2. Shields.io Badges
```
https://img.shields.io/badge/{label}-{message}-{color}?logo={logo}
```

**Примеры:**
```
https://img.shields.io/badge/Go-1.25-00ADD8?logo=go
https://img.shields.io/badge/License-MIT-green.svg
```

**Генератор:** https://shields.io/

#### 3. Simple Icons
```
https://cdn.simpleicons.org/{name}/{color}
```

---

## 📱 Mobile Responsive

### Адаптивная сетка

Автоматически адаптируется:

**Desktop (> 1024px):**
```
┌─────────┬──────────────────────┐
│ Sidebar │ Контент (3 колонки)  │
└─────────┴──────────────────────┘
```

**Tablet (768px - 1024px):**
```
┌─────────┬──────────────────────┐
│ Sidebar │ Контент (2 колонки)  │
└─────────┴──────────────────────┘
```

**Mobile (< 768px):**
```
┌──────────────────────┐
│ Hamburger Menu ☰    │
├──────────────────────┤
│ Контент (1 колонка) │
└──────────────────────┘
```

---

## 🚀 GitHub Pages Setup

**Полная инструкция:** [GitHub Pages Setup Guide](./GITHUB_PAGES_SETUP.md)

**Кратко:**

1. Push на GitHub
2. Settings → Pages → Source: main, /docs
3. Сайт доступен: `https://alexrodving.github.io/advanced-user-api/`

---

## 🔧 Как обновить сайт

### Изменить контент

```bash
# Отредактируйте Markdown файлы в docs/docs/
vim docs/docs/API.md

# Коммит и push
git add docs/
git commit -m "docs: update API documentation"
git push origin main

# GitHub Pages обновится автоматически через 1-2 минуты
```

### Изменить дизайн

```bash
# Отредактируйте CSS в docs/index.html
vim docs/index.html

# Коммит и push
git add docs/index.html
git commit -m "style: update website design"
git push origin main
```

### Добавить новую страницу

```bash
# Создайте новый MD файл
echo "# New Page" > docs/docs/NEW_PAGE.md

# Добавьте в sidebar
vim docs/_sidebar.md
# Добавьте: - [New Page](docs/NEW_PAGE.md)

# Коммит
git add docs/
git commit -m "docs: add new page"
git push origin main
```

---

## 🎯 Best Practices

### 1. Оптимизация изображений

```html
<!-- Используйте CDN для иконок -->
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/...">

<!-- Или оптимизированные SVG -->
<img src="images/logo.svg" width="60" height="60">
```

### 2. Внутренние ссылки

```markdown
<!-- Относительные пути -->
[API Docs](docs/API.md)

<!-- Якоря -->
[Перейти к разделу](#section-name)
```

### 3. Кеширование

CDN библиотеки автоматически кешируются браузером

### 4. SEO

```html
<meta name="description" content="Production-ready REST API with JWT">
<meta property="og:title" content="Advanced User API">
<meta property="og:description" content="...">
```

---

## 🆚 Почему Docsify а не другие?

### Docsify ✅
- Не требует сборки
- Работает прямо с Markdown
- Быстрая загрузка
- Множество плагинов
- Легко настраивается

### MkDocs
- Требует Python
- Нужна сборка
- Больше настроек
- Material theme красивый

### Docusaurus
- Требует Node.js
- React-based
- Тяжёлая сборка
- Мощный но избыточный

### VitePress
- Требует Node.js
- Vue-based
- Быстрая сборка
- Хорош для Vue проектов

**Для Go проекта Docsify идеален!**

---

## 📖 Дополнительные ресурсы

### Официальные документации
- [Docsify Documentation](https://docsify.js.org/)
- [GitHub Pages Guide](https://pages.github.com/)
- [Markdown Guide](https://www.markdownguide.org/)

### Инструменты дизайна
- [UI Gradients](https://uigradients.com/) - красивые градиенты
- [DevIcons](https://devicon.dev/) - иконки технологий
- [Shields.io](https://shields.io/) - badges
- [CSS Gradient Generator](https://cssgradient.io/)

### CSS References
- [MDN CSS Grid](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Grid_Layout)
- [MDN Flexbox](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout)
- [CSS Tricks](https://css-tricks.com/)

---

## 🎨 Кастомизация

### Изменить цветовую схему

```css
/* docs/index.html */
:root {
  --theme-color: #YOUR_COLOR;
}

.hero {
  background: linear-gradient(135deg, #COLOR1 0%, #COLOR2 100%);
}
```

### Добавить Google Analytics

```html
<!-- В <head> секции docs/index.html -->
<script async src="https://www.googletagmanager.com/gtag/js?id=GA_ID"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'GA_ID');
</script>
```

### Сменить иконки

Замените URL в `docs/README.md`:
```html
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg">
```

Найдите другие на: https://devicon.dev/

---

## 🐛 Troubleshooting

### Сайт показывает 404

**Проблема:** GitHub Pages не настроен

**Решение:**
1. Settings → Pages
2. Source: main, /docs
3. Подождите 5 минут

---

### Стили не применяются

**Проблема:** Файлы с `_` игнорируются Jekyll

**Решение:** Убедитесь что есть `.nojekyll` файл в `docs/`

---

### Sidebar не показывается

**Проблема:** Неправильный путь в конфигурации

**Решение:** Проверьте `index.html`:
```javascript
loadSidebar: true,  // Должно быть true
```

---

### Иконки не загружаются

**Проблема:** CDN недоступен

**Решение:** Используйте другой CDN или локальные файлы

---

## 📊 Производительность

### Оптимизации Docsify:

1. **Lazy loading** - страницы загружаются по требованию
2. **CDN caching** - библиотеки кешируются
3. **Минификация** - используются minified версии
4. **No build step** - мгновенное обновление

### Метрики:

- **First Load:** < 1s
- **Page Switch:** < 100ms
- **Search:** < 50ms
- **Total Size:** < 500KB (без документации)

---

## 🎓 Что вы узнали

### Frontend:
- ✅ HTML структура
- ✅ CSS Grid и Flexbox
- ✅ CSS transitions и animations
- ✅ Responsive design
- ✅ CDN использование

### Docsify:
- ✅ Конфигурация
- ✅ Плагины
- ✅ Sidebar/Navbar
- ✅ Темизация

### GitHub Pages:
- ✅ Настройка Pages
- ✅ Jekyll vs no-Jekyll
- ✅ Деплой статических сайтов

---

**Ваш профессиональный сайт документации готов!** 🎉

