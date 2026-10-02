const supportedSignatures: Array<{ mime: string; matches: (bytes: Uint8Array) => boolean }> = [
  { mime: 'image/png', matches: (b) => b[0] === 0x89 && ascii(b, 1, 4) === 'PNG' },
  { mime: 'image/jpeg', matches: (b) => b[0] === 0xff && b[1] === 0xd8 && b[2] === 0xff },
  { mime: 'image/gif', matches: (b) => ascii(b, 0, 4) === 'GIF8' },
  { mime: 'image/webp', matches: (b) => ascii(b, 0, 4) === 'RIFF' && ascii(b, 8, 12) === 'WEBP' },
  { mime: 'image/bmp', matches: (b) => ascii(b, 0, 2) === 'BM' },
  { mime: 'image/x-icon', matches: (b) => b[0] === 0 && b[1] === 0 && b[2] === 1 && b[3] === 0 },
]

const maxImageBytes = 20 * 1024 * 1024

export function normalizeImageDataUrl(value: string): string {
  const input = value.trim()
  if (!input) throw new Error('请粘贴图片 Base64 内容')

  const dataUrl = input.match(/^data:(image\/[\w.+-]+);base64,([\s\S]+)$/i)
  const declaredMime = dataUrl?.[1].toLowerCase()
  const payload = (dataUrl?.[2] ?? input).replace(/\s/g, '')
  if (!payload || !/^[A-Za-z0-9+/]*={0,2}$/.test(payload)) {
    throw new Error('Base64 内容包含无效字符')
  }
  if (Math.ceil(payload.length * 0.75) > maxImageBytes) {
    throw new Error('图片不能超过 20 MB')
  }

  let binary: string
  try {
    binary = atob(payload)
  } catch {
    throw new Error('无法解析 Base64 内容')
  }
  const bytes = Uint8Array.from(binary, (char) => char.charCodeAt(0))
  const detectedMime = supportedSignatures.find((signature) => signature.matches(bytes))?.mime
  if (!detectedMime) {
    throw new Error('无法识别图片格式，请使用 PNG、JPEG、GIF、WebP、BMP 或 ICO')
  }
  if (declaredMime && normalizeMime(declaredMime) !== detectedMime) {
    throw new Error('图片声明格式与实际内容不一致')
  }
  return `data:${detectedMime};base64,${payload}`
}

function ascii(bytes: Uint8Array, start: number, end: number): string {
  return String.fromCharCode(...bytes.slice(start, end))
}

function normalizeMime(mime: string): string {
  return mime === 'image/vnd.microsoft.icon' ? 'image/x-icon' : mime
}
