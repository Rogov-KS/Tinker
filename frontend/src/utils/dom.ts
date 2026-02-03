export function createElement<K extends keyof HTMLElementTagNameMap>(
  tag: K,
  options: {
    className?: string
    text?: string
    html?: string
    attrs?: Record<string, string>
  } = {}
): HTMLElementTagNameMap[K] {
  const element = document.createElement(tag)
  if (options.className) element.className = options.className
  if (options.text) element.textContent = options.text
  if (options.html) element.innerHTML = options.html
  if (options.attrs) {
    Object.entries(options.attrs).forEach(([key, value]) =>
      element.setAttribute(key, value)
    )
  }
  return element
}

export function clearChildren(node: HTMLElement) {
  while (node.firstChild) {
    node.removeChild(node.firstChild)
  }
}

/**
 * Создает элемент аватара с поддержкой изображения из URL или fallback на первую букву имени
 * @param username - имя пользователя для fallback
 * @param avatarUrl - опциональный URL изображения аватара
 * @param size - размер аватара: 'small' | 'large' | undefined (обычный)
 * @returns HTMLDivElement с классом 'avatar'
 */
export function createAvatar(
  username: string,
  avatarUrl?: string | null,
  size?: 'small' | 'large'
): HTMLDivElement {
  const avatar = createElement('div', {
    className: `avatar${size ? ` ${size}` : ''}`,
  })

  // Если есть URL, пытаемся загрузить изображение
  if (avatarUrl && avatarUrl.trim()) {
    const img = document.createElement('img')
    img.src = avatarUrl
    img.alt = username
    img.style.width = '100%'
    img.style.height = '100%'
    img.style.objectFit = 'cover'
    img.style.borderRadius = '50%'
    img.style.position = 'absolute'
    img.style.top = '0'
    img.style.left = '0'

    // Создаем элемент для текста (fallback)
    const textFallback = document.createElement('span')
    textFallback.textContent = (username[0] || 'T').toUpperCase()
    textFallback.style.position = 'relative'
    textFallback.style.zIndex = '1'

    // При ошибке загрузки показываем букву, скрываем изображение
    img.addEventListener('error', () => {
      img.style.display = 'none'
      textFallback.style.display = 'block'
    })

    // При успешной загрузке скрываем текст
    img.addEventListener('load', () => {
      textFallback.style.display = 'none'
    })

    avatar.appendChild(textFallback)
    avatar.appendChild(img)
  } else {
    // Если нет URL, сразу показываем букву
    avatar.textContent = (username[0] || 'T').toUpperCase()
  }

  return avatar
}


