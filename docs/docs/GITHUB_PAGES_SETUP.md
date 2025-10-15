# 🌐 GitHub Setup & Pages

## 📋 Шаг 1: Создайте репозиторий на GitHub

1. Откройте: https://github.com/new

2. Заполните:
   - **Repository name:** `advanced-user-api`
   - **Description:** `Production-ready REST API with JWT authentication, Gin framework, GORM ORM, Docker, and CI/CD`
   - **Visibility:** Public ✅ (для портфолио)
   
3. **НЕ добавляйте:**
   - ❌ README
   - ❌ .gitignore
   - ❌ License
   
   (у вас уже всё есть)

4. Нажмите **"Create repository"**

5. **Topics (теги):** Добавьте на странице репозитория (⚙️ рядом с About):
   ```
   golang, gin, gorm, jwt, docker, rest-api, postgresql, clean-architecture, ci-cd
   ```

---

## 🔗 Шаг 2: Добавьте remote и запуште

```bash
# Перейдите в директорию проекта
cd advanced-user-api

# Добавьте remote (SSH - рекомендуется)
git remote add origin git@github.com:AlexRodving/advanced-user-api.git

# Или HTTPS (если нет SSH ключа)
git remote add origin https://github.com/AlexRodving/advanced-user-api.git

# Запуште код
git push -u origin main
```

**Если используете HTTPS и нужна аутентификация:**
- Используйте Personal Access Token вместо пароля
- GitHub → Settings → Developer settings → Personal access tokens → Generate new token

---

## 🌐 Шаг 3: Настройте GitHub Pages

1. Откройте ваш репозиторий на GitHub:
   ```
   https://github.com/AlexRodving/advanced-user-api
   ```

2. Перейдите в **Settings** → **Pages** (в левом меню)

3. В разделе **Build and deployment**:
   - **Source**: Deploy from a branch
   - **Branch**: выберите `main`
   - **Folder**: выберите `/docs`
   - Нажмите **Save**

4. Готово! Через 2-3 минуты сайт будет доступен:
   ```
   https://alexrodving.github.io/advanced-user-api/
   ```

---

## ✅ Шаг 4: Проверка

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
# Отредактируйте файлы в docs/docs/
vim docs/docs/API.md

# Коммит и пуш
git add docs/
git commit -m "docs: update API documentation"
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

