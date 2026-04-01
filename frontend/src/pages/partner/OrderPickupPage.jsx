import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { observer } from 'mobx-react-lite'
import { BadgeCheck, Camera, Clock3, Keyboard, PackageCheck, QrCode, StopCircle, User } from 'lucide-react'
import { Html5Qrcode, Html5QrcodeSupportedFormats } from 'html5-qrcode'
import { toast } from 'sonner'
import { formatDateTime, getErrorMessage } from '@/lib/utils'
import PartnerLayout from '@/components/partner/layout/PartnerLayout'
import Input from '@/components/ui/form/Input'
import Button from '@/components/ui/actions/Button'
import { useStores } from '@/context/StoresContext'

const SCANNER_ID = 'partner-order-qr-reader'

function OrderPickupPageBase() {
  const { ordersStore } = useStores()
  const scannerRef = useRef(null)
  const isMountedRef = useRef(true)
  const hasScannedRef = useRef(false)

  const [pickupCode, setPickupCode] = useState('')
  const [isScannerActive, setIsScannerActive] = useState(false)
  const [scannerError, setScannerError] = useState('')

  const order = ordersStore.current

  const stopScanner = useCallback(async () => {
    const scanner = scannerRef.current
    if (!scanner) {
      setIsScannerActive(false)
      return
    }

    try {
      if (scanner.isScanning) await scanner.stop()
    } catch {
      // noop
    }

    try {
      await scanner.clear()
    } catch {
      // noop
    }

    scannerRef.current = null
    hasScannedRef.current = false
    if (isMountedRef.current) setIsScannerActive(false)
  }, [])

  const handleLookup = useCallback(async (rawCode) => {
    const code = (rawCode || '').trim().toUpperCase()
    if (!code) {
      toast.error('Введите код получения')
      return
    }

    setPickupCode(code)
    try {
      await ordersStore.lookupByCode(code)
      toast.success('Заказ найден')
    } catch (error) {
      if (error?.response?.status === 404) {
        toast.error('Заказ с таким кодом не найден')
        return
      }
      toast.error(getErrorMessage(error))
    }
  }, [ordersStore])

  const handleScanSuccess = useCallback((decodedText) => {
    if (hasScannedRef.current) return
    const code = (decodedText || '').trim().toUpperCase()
    if (!code) return

    hasScannedRef.current = true
    stopScanner()
    handleLookup(code)
  }, [handleLookup, stopScanner])

  const startScanner = useCallback(async () => {
    if (isScannerActive) return
    setScannerError('')

    try {
      const cameras = await Html5Qrcode.getCameras()
      if (!cameras?.length) {
        setScannerError('Камера не найдена на устройстве')
        return
      }

      const backCamera = cameras.find((camera) => /back|rear|environment|traseira|trasera/i.test(camera.label))
      const cameraConfig = backCamera ? { deviceId: { exact: backCamera.id } } : { facingMode: 'environment' }

      const scanner = new Html5Qrcode(SCANNER_ID)
      scannerRef.current = scanner

      await scanner.start(
        cameraConfig,
        { fps: 10, qrbox: { width: 240, height: 240 }, formatsToSupport: [Html5QrcodeSupportedFormats.QR_CODE] },
        handleScanSuccess,
        () => {}
      )

      setIsScannerActive(true)
    } catch (error) {
      try {
        if (scannerRef.current) await scannerRef.current.clear()
      } catch {
        // noop
      }
      scannerRef.current = null
      setIsScannerActive(false)
      setScannerError(error?.message || 'Не удалось запустить камеру')
    }
  }, [handleScanSuccess, isScannerActive])

  useEffect(() => {
    return () => {
      isMountedRef.current = false
      stopScanner()
    }
  }, [stopScanner])

  const statusBadge = useMemo(() => {
    return STATUS_META[order?.status] || { label: order?.status || '—', className: 'bg-gray-100 text-gray-800' }
  }, [order?.status])

  const canIssueOrder = order?.status === 'confirmed'

  const handlePickup = async () => {
    try {
      const data = await ordersStore.pickup(order.id)
      toast.success(data?.message || 'Заказ отмечен как выданный')
    } catch (error) {
      if (error?.response?.status === 409) {
        toast.error('Заказ нельзя выдать в текущем статусе')
        return
      }
      toast.error(getErrorMessage(error))
    }
  }

  return (
    <PartnerLayout title="Выдача заказа" subtitle="Отсканируйте QR или введите код получения вручную">
      <div className="space-y-6 max-w-6xl">
        <div className="grid lg:grid-cols-2 gap-5">
          <section className="card space-y-4">
            <h2 className="text-base font-semibold text-brand-900 flex items-center gap-2">
              <Camera size={18} className="text-brand-500" /> Сканирование QR
            </h2>
            <div className="rounded-xl border border-cream-200 bg-white p-4">
              <div id={SCANNER_ID} className="min-h-[280px] w-full overflow-hidden rounded-lg bg-cream-100" />
            </div>
            {scannerError && <p className="text-sm text-red-600">{scannerError}</p>}
            <div className="flex gap-3 flex-wrap">
              <Button onClick={startScanner} disabled={isScannerActive} className="gap-2">
                <QrCode size={16} /> Запустить сканер
              </Button>
              <Button variant="secondary" onClick={stopScanner} disabled={!isScannerActive} className="gap-2">
                <StopCircle size={16} /> Остановить
              </Button>
            </div>
          </section>

          <section className="card space-y-4">
            <h2 className="text-base font-semibold text-brand-900 flex items-center gap-2">
              <Keyboard size={18} className="text-brand-500" /> Ручной ввод кода
            </h2>
            <Input
              value={pickupCode}
              onChange={(e) => setPickupCode(e.target.value.toUpperCase())}
              placeholder="Например, AB12CD34"
              autoCapitalize="characters"
              autoCorrect="off"
              spellCheck={false}
              maxLength={32}
            />
            <Button onClick={() => handleLookup(pickupCode)} disabled={ordersStore.lookupLoading || !pickupCode.trim()} className="w-full">
              {ordersStore.lookupLoading ? 'Ищем заказ...' : 'Найти заказ'}
            </Button>
          </section>
        </div>

        <section className="card space-y-4">
          <h2 className="text-base font-semibold text-brand-900 flex items-center gap-2">
            <PackageCheck size={18} className="text-brand-500" /> Информация о заказе
          </h2>

          {!order && (
            <div className="rounded-xl border-2 border-dashed border-cream-300 p-8 text-center text-brand-500">
              Сначала найдите заказ по QR или коду получения
            </div>
          )}

          {order && (
            <div className="space-y-5">
              <div className="flex items-start justify-between flex-wrap gap-3">
                <div>
                  <p className="text-sm text-brand-500">Код получения</p>
                  <p className="text-lg font-semibold tracking-wide">{order.pickup_code}</p>
                </div>
                <span className={`badge ${statusBadge.className}`}>{statusBadge.label}</span>
              </div>

              <div className="grid md:grid-cols-2 gap-4">
                <InfoRow label="Бокс" value={order.box?.name || '—'} icon={PackageCheck} />
                <InfoRow label="Клиент" value={order.customer?.name ? `${order.customer.name} (${order.customer.phone})` : order.customer?.phone || '—'} icon={User} />
                <InfoRow label="Окно выдачи" value={`${formatDateTime(order.pickup_time?.start)} - ${formatDateTime(order.pickup_time?.end)}`} icon={Clock3} />
                <InfoRow label="Создан" value={formatDateTime(order.created_at)} icon={BadgeCheck} />
              </div>

              <Button onClick={handlePickup} disabled={!canIssueOrder || ordersStore.pickupLoading} className="w-full sm:w-auto gap-2">
                <PackageCheck size={16} />
                {ordersStore.pickupLoading ? 'Выдаем заказ...' : 'Выдать заказ'}
              </Button>
            </div>
          )}
        </section>
      </div>
    </PartnerLayout>
  )
}

const STATUS_META = {
  confirmed: { label: 'Подтвержден', className: 'bg-green-100 text-green-800' },
  picked_up: { label: 'Выдан', className: 'bg-blue-100 text-blue-800' },
  completed: { label: 'Завершен', className: 'bg-brand-100 text-brand-800' },
  pending: { label: 'Ожидает', className: 'bg-yellow-100 text-yellow-800' },
  paid: { label: 'Оплачен', className: 'bg-emerald-100 text-emerald-800' },
  cancelled: { label: 'Отменен', className: 'bg-red-100 text-red-800' },
}

function InfoRow({ label, value, icon: Icon }) {
  return (
    <div className="rounded-xl border border-cream-200 p-3 bg-cream-50">
      <p className="text-xs uppercase tracking-wider text-cream-500 mb-1">{label}</p>
      <p className="text-sm font-medium text-brand-800 flex items-center gap-1.5">
        {Icon && <Icon size={14} className="text-brand-400" />}
        {value || '—'}
      </p>
    </div>
  )
}

export default observer(OrderPickupPageBase)
