import { useRef, useState } from 'react'
import { Upload, X } from 'lucide-react'
import { toast } from 'sonner'
import { cn } from '@/lib/utils'
import { uploadMedia } from '@/api/partner'
import Spinner from '@/components/ui/feedback/Spinner'

export default function ImageUpload({ value, onChange, error }) {
  const [isUploading, setIsUploading] = useState(false)
  const [preview, setPreview] = useState(value || null)
  const inputRef = useRef(null)

  const handleFile = async (file) => {
    if (!file) return
    if (!file.type.startsWith('image/')) {
      toast.error('Выберите изображение')
      return
    }
    if (file.size > 5 * 1024 * 1024) {
      toast.error('Размер файла должен быть до 5 МБ')
      return
    }

    const reader = new FileReader()
    reader.onloadend = () => setPreview(reader.result)
    reader.readAsDataURL(file)

    setIsUploading(true)
    try {
      const uploaded = await uploadMedia(file)
      onChange(uploaded.url)
      toast.success('Изображение загружено')
    } catch {
      toast.error('Не удалось загрузить изображение')
      setPreview(null)
    } finally {
      setIsUploading(false)
    }
  }

  const handleRemove = () => {
    setPreview(null)
    onChange('')
    if (inputRef.current) inputRef.current.value = ''
  }

  return (
    <div>
      <label className="block text-sm font-medium text-brand-700 mb-2">Изображение бокса</label>
      {preview ? (
        <div className="relative w-full h-64 rounded-xl overflow-hidden border border-cream-200">
          <img src={preview} alt="preview" className="w-full h-full object-cover" />
          {isUploading && (
            <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
              <Spinner size={28} className="text-white" />
            </div>
          )}
          {!isUploading && (
            <button type="button" onClick={handleRemove} className="absolute top-2 right-2 p-2 bg-red-500 text-white rounded-lg hover:bg-red-600">
              <X size={16} />
            </button>
          )}
        </div>
      ) : (
        <button
          type="button"
          onClick={() => inputRef.current?.click()}
          className={cn(
            'w-full h-56 rounded-xl border-2 border-dashed bg-cream-50 flex flex-col items-center justify-center gap-2 transition-colors',
            error ? 'border-red-400' : 'border-cream-300 hover:border-brand-400'
          )}
        >
          {isUploading ? <Spinner size={28} /> : <Upload size={26} className="text-brand-500" />}
          <span className="text-sm text-brand-700">Нажмите, чтобы загрузить</span>
          <span className="text-xs text-brand-500">PNG, JPG до 5 МБ</span>
        </button>
      )}

      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={(e) => handleFile(e.target.files?.[0])}
      />

      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  )
}
