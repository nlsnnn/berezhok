import { useState, useRef } from 'react'
import { Upload, X, ImageIcon } from 'lucide-react'
import { cn } from '@/lib/utils'
import { uploadMedia } from '@/api/partner'
import { toast } from 'sonner'
import Spinner from './Spinner'

export default function ImageUpload({ value, onChange, error }) {
  const [isUploading, setIsUploading] = useState(false)
  const [preview, setPreview] = useState(value || null)
  const fileInputRef = useRef(null)

  const handleFileChange = async (e) => {
    const file = e.target.files?.[0]
    if (!file) return

    // Validate file type
    if (!file.type.startsWith('image/')) {
      toast.error('Пожалуйста, выберите изображение')
      return
    }

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      toast.error('Размер файла не должен превышать 5 МБ')
      return
    }

    // Create preview
    const reader = new FileReader()
    reader.onloadend = () => {
      setPreview(reader.result)
    }
    reader.readAsDataURL(file)

    // Upload file
    setIsUploading(true)
    try {
      const response = await uploadMedia(file)
      onChange(response.url)
      toast.success('Изображение загружено')
    } catch (err) {
      console.error('Upload error:', err)
      toast.error('Не удалось загрузить изображение')
      setPreview(null)
    } finally {
      setIsUploading(false)
    }
  }

  const handleRemove = () => {
    setPreview(null)
    onChange('')
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const handleDragOver = (e) => {
    e.preventDefault()
  }

  const handleDrop = (e) => {
    e.preventDefault()
    const file = e.dataTransfer.files?.[0]
    if (file) {
      fileInputRef.current.files = e.dataTransfer.files
      handleFileChange({ target: { files: [file] } })
    }
  }

  return (
    <div className="w-full">
      <label className="block text-sm font-medium text-brand-700 mb-2">
        Изображение бокса
      </label>

      {preview ? (
        <div className="relative w-full h-64 rounded-xl overflow-hidden border-2 border-cream-200 group">
          <img
            src={preview}
            alt="Preview"
            className="w-full h-full object-cover"
          />
          {isUploading && (
            <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center">
              <Spinner size={32} className="text-white" />
            </div>
          )}
          {!isUploading && (
            <button
              type="button"
              onClick={handleRemove}
              className="absolute top-2 right-2 p-2 bg-red-500 text-white rounded-lg opacity-0 group-hover:opacity-100 transition-opacity hover:bg-red-600"
            >
              <X size={18} />
            </button>
          )}
        </div>
      ) : (
        <div
          onDragOver={handleDragOver}
          onDrop={handleDrop}
          className={cn(
            'w-full h-64 rounded-xl border-2 border-dashed flex flex-col items-center justify-center cursor-pointer transition-colors',
            error ? 'border-red-400 bg-red-50' : 'border-cream-300 bg-cream-50 hover:border-brand-400 hover:bg-brand-50',
            isUploading && 'pointer-events-none opacity-50'
          )}
          onClick={() => fileInputRef.current?.click()}
        >
          {isUploading ? (
            <Spinner size={32} />
          ) : (
            <>
              <div className="w-16 h-16 rounded-full bg-brand-100 flex items-center justify-center mb-4">
                <Upload size={28} className="text-brand-500" />
              </div>
              <p className="text-sm font-medium text-brand-800 mb-1">
                Нажмите или перетащите изображение
              </p>
              <p className="text-xs text-brand-500">
                PNG, JPG до 5 МБ
              </p>
            </>
          )}
        </div>
      )}

      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileChange}
        className="hidden"
      />

      {error && <p className="mt-2 text-xs text-red-500">{error}</p>}
    </div>
  )
}
