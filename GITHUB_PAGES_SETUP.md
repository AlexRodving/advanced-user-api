# 🌐 Настройка GitHub Pages

## Шаг 1: Подготовка

Весь сайт уже готов в папке `docs-site/`!

```
docs-site/
├── index.html          # Главная страница Docsify
├── README.md           # Домашняя страница сайта
├── _sidebar.md         # Боковое меню
├── _navbar.md          # Верхнее меню
└── docs/               # Вся документация
    ├── API.md
    ├── ARCHITECTURE.md
    ├── TESTING.md
    └── libraries/
        ├── GIN.md
        ├── GORM.md
        ├── JWT.md
        └── ...
```

---

## Шаг 2: Пуш на GitHub

```bash
cd /home/rodving/Документы/go/teach/08_advanced_api

# Добавляем docs-site
git add docs-site/
git commit -m "feat: add GitHub Pages documentation site

- Beautiful Docsify-based documentation site
- Custom styling with gradients and cards
- Tech stack icons from devicons
- Full documentation integrated
- Sidebar navigation
- Search functionality
- Syntax highlighting for Go, SQL, Bash
- Mobile responsive design"

# Пушим на GitHub
git push origin main
```

---

## Шаг 3: Настройка GitHub Pages

### Через GitHub UI:

1. Откройте ваш репозиторий на GitHub:
   ```
   https://github.com/AlexRodving/advanced-user-api
   ```

2. Перейдите в **Settings** (вверху справа)

3. В левом меню выберите **Pages**

4. В разделе **Source**:
   - **Branch**: выберите `main`
   - **Folder**: выберите `/docs-site`
   - Нажмите **Save**

5. Готово! Через 1-2 минуты сайт будет доступен по адресу:
   ```
   https://alexrodving.github.io/advanced-user-api/
   ```

---

## Шаг 4: Проверка

Откройте в браузере:
```
https://alexrodving.github.io/advanced-user-api/
```

Вы увидите:
- 🎨 Красивую главную страницу с градиентами
- 🖼️ Иконки технологий (Go, PostgreSQL, Docker, Redis)
- 📊 Статистику проекта
- 📚 Полную документацию в боковом меню
- 🔍 Поиск по документации
- 📱 Адаптивный дизайн для мобильных

---

## Особенности сайта

### Дизайн
- ✅ Gradient hero section
- ✅ Tech stack icons from devicons
- ✅ Feature cards с hover эффектами
- ✅ Statistics counters
- ✅ Alert boxes (info, success, warning)
- ✅ Custom CSS styling
- ✅ Mobile responsive

### Функционал
- ✅ **Search** - поиск по всей документации
- ✅ **Syntax highlighting** - Go, SQL, Bash, JSON, Docker
- ✅ **Copy code** - копирование кода одним кликом
- ✅ **Pagination** - переход между страницами
- ✅ **Emoji** - поддержка эмодзи
- ✅ **Sidebar navigation** - удобное меню
- ✅ **Auto scroll to top** - автоскролл при переходе

---

## Обновление сайта

При изменении документации:

```bash
# Обновите docs-site/docs/
# (они автоматически копируются из docs/)

# Или скопируйте изменения:
cp -r docs/* docs-site/docs/

# Коммит и пуш
git add docs-site/
git commit -m "docs: update documentation"
git push origin main

# GitHub Pages обновится автоматически через 1-2 минуты
```

---

## Кастомизация

### Изменить цветовую тему

Отредактируйте `docs-site/index.html`:

```css
:root {
  --theme-color: #00ADD8;  /* Основной цвет (сейчас Go голубой) */
  --theme-color-dark: #007d9c;
}
```

### Изменить градиент hero

```css
.hero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  /* Попробуйте другие градиенты с uigradients.com */
}
```

### Добавить Google Analytics

Добавьте в `docs-site/index.html` перед `</head>`:

```html
<!-- Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id=GA_TRACKING_ID"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'GA_TRACKING_ID');
</script>
```

---

## Возможные проблемы

### Сайт не появляется

1. Проверьте что `docs-site/` папка в main ветке
2. Проверьте Settings → Pages → Source = `main` branch, `/docs-site` folder
3. Подождите 5-10 минут (первая публикация может занять время)

### 404 ошибки на страницах

Проверьте что все ссылки в документации относительные:
```markdown
<!-- ✅ Правильно -->
[API](docs/API.md)

<!-- ❌ Неправильно -->
[API](/docs/API.md)
```

### Стили не применяются

Проверьте что `index.html` в корне `docs-site/`

---

## Альтернативы Docsify

Если захотите другой генератор:

### MkDocs (Python)
```bash
pip install mkdocs mkdocs-material
mkdocs new .
mkdocs serve
```

### Docusaurus (React)
```bash
npx create-docusaurus@latest my-website classic
npm start
```

### VitePress (Vue)
```bash
npm init vitepress
npm run docs:dev
```

Но **Docsify** самый простой - не требует сборки, работает прямо с Markdown!

---

## Полезные ссылки

- [Docsify Documentation](https://docsify.js.org/)
- [GitHub Pages Guide](https://pages.github.com/)
- [DevIcons](https://devicon.dev/) - иконки технологий
- [UI Gradients](https://uigradients.com/) - красивые градиенты

---

**Ваш сайт документации готов!** 🎉

Откройте: `https://alexrodving.github.io/advanced-user-api/`

