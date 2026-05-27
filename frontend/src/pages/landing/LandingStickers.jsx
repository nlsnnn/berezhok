const STK_STROKE = "#3D2818"

export function StickerBread({ size = 80, tilt = -8 }) {
  return (
    <svg viewBox="0 0 100 100" width={size} height={size} style={{ transform: `rotate(${tilt}deg)` }}>
      <ellipse cx="50" cy="50" rx="42" ry="32" fill="#fff" stroke="#fff" strokeWidth="6" />
      <ellipse cx="50" cy="52" rx="38" ry="28" fill="#E0A86A" />
      <ellipse cx="50" cy="50" rx="38" ry="28" fill="#D49858" />
      <path d="M 26 52 Q 50 38 74 52" stroke="#8B5A2B" strokeWidth="2.5" fill="none" opacity="0.5" />
      <path d="M 32 44 L 38 50 M 44 40 L 50 46 M 56 40 L 62 46 M 68 44 L 74 50" stroke="#8B5A2B" strokeWidth="2.5" strokeLinecap="round" />
      <ellipse cx="38" cy="46" rx="3" ry="1.5" fill="#fff" opacity="0.5" />
    </svg>
  )
}

export function StickerCroissant({ size = 80, tilt = 10 }) {
  return (
    <svg viewBox="0 0 100 100" width={size} height={size} style={{ transform: `rotate(${tilt}deg)` }}>
      <path d="M 22 40 Q 30 18 50 16 Q 76 14 82 38 Q 86 60 70 70 Q 50 78 38 70 Q 18 60 22 40 Z"
        fill="#fff" stroke="#fff" strokeWidth="6" strokeLinejoin="round" />
      <path d="M 22 40 Q 30 18 50 16 Q 76 14 82 38 Q 86 60 70 70 Q 50 78 38 70 Q 18 60 22 40 Z"
        fill="#E8A85A" stroke={STK_STROKE} strokeWidth="2" strokeLinejoin="round" />
      <path d="M 32 38 L 40 28 M 44 36 L 52 24 M 58 34 L 64 24 M 36 50 L 46 44 M 54 50 L 64 42 M 38 62 L 50 58 M 56 62 L 68 56"
        stroke="#8B4513" strokeWidth="2" strokeLinecap="round" opacity="0.55" />
    </svg>
  )
}

export function StickerCoffee({ size = 70, tilt = -6 }) {
  return (
    <svg viewBox="0 0 100 100" width={size} height={size} style={{ transform: `rotate(${tilt}deg)` }}>
      <path d="M 24 30 L 30 80 Q 32 88 50 88 Q 68 88 70 80 L 76 30 Z" fill="#fff" stroke="#fff" strokeWidth="6" strokeLinejoin="round" />
      <path d="M 24 30 L 30 80 Q 32 88 50 88 Q 68 88 70 80 L 76 30 Z" fill="#F4DCC4" stroke={STK_STROKE} strokeWidth="2.5" strokeLinejoin="round" />
      <rect x="22" y="22" width="56" height="14" rx="4" fill="#A85E2C" stroke={STK_STROKE} strokeWidth="2.5" />
      <circle cx="50" cy="29" r="3" fill="#3D2818" />
      <path d="M 38 16 Q 42 8 38 0 M 50 14 Q 54 6 50 -2 M 62 16 Q 66 8 62 0" stroke={STK_STROKE} strokeWidth="2.5" fill="none" strokeLinecap="round" opacity="0.55" transform="translate(0,4)" />
      <rect x="28" y="50" width="44" height="14" fill="#D89A52" stroke={STK_STROKE} strokeWidth="2" />
    </svg>
  )
}

export function StickerCoin({ size = 70, tilt = 4 }) {
  return (
    <svg viewBox="0 0 100 100" width={size} height={size} style={{ transform: `rotate(${tilt}deg)` }}>
      <circle cx="50" cy="50" r="40" fill="#fff" stroke="#fff" strokeWidth="6" />
      <circle cx="50" cy="50" r="36" fill="#F4C04B" stroke={STK_STROKE} strokeWidth="2.5" />
      <circle cx="50" cy="50" r="28" fill="none" stroke="#B88C24" strokeWidth="2" opacity="0.6" />
      <text x="50" y="64" textAnchor="middle" fill="#8C6A1B" fontFamily="Unbounded, system-ui" fontWeight="800" fontSize="38">₽</text>
      <ellipse cx="38" cy="38" rx="8" ry="4" fill="#fff" opacity="0.4" />
    </svg>
  )
}

export function StickerHeart({ size = 60, tilt = -10 }) {
  return (
    <svg viewBox="0 0 100 100" width={size} height={size} style={{ transform: `rotate(${tilt}deg)` }}>
      <path d="M 50 84 C 18 60 14 30 32 22 C 42 18 50 26 50 32 C 50 26 58 18 68 22 C 86 30 82 60 50 84 Z"
        fill="#fff" stroke="#fff" strokeWidth="6" strokeLinejoin="round" />
      <path d="M 50 84 C 18 60 14 30 32 22 C 42 18 50 26 50 32 C 50 26 58 18 68 22 C 86 30 82 60 50 84 Z"
        fill="#E85A6B" stroke={STK_STROKE} strokeWidth="2.5" strokeLinejoin="round" />
      <ellipse cx="38" cy="40" rx="6" ry="3" fill="#fff" opacity="0.4" />
    </svg>
  )
}

export function StickerSparkle({ size = 50, tilt = 0 }) {
  return (
    <svg viewBox="0 0 100 100" width={size} height={size} style={{ transform: `rotate(${tilt}deg)` }}>
      <path d="M 50 5 L 56 44 L 95 50 L 56 56 L 50 95 L 44 56 L 5 50 L 44 44 Z"
        fill="#F4C04B" stroke={STK_STROKE} strokeWidth="2" />
    </svg>
  )
}

export function StickerLeaf({ size = 60, tilt = 14 }) {
  return (
    <svg viewBox="0 0 100 100" width={size} height={size} style={{ transform: `rotate(${tilt}deg)` }}>
      <path d="M 18 80 Q 14 30 50 14 Q 88 22 82 60 Q 76 88 50 88 Q 28 88 18 80 Z"
        fill="#fff" stroke="#fff" strokeWidth="6" strokeLinejoin="round" />
      <path d="M 18 80 Q 14 30 50 14 Q 88 22 82 60 Q 76 88 50 88 Q 28 88 18 80 Z"
        fill="#6FAA56" stroke={STK_STROKE} strokeWidth="2.5" strokeLinejoin="round" />
      <path d="M 22 78 Q 48 50 80 30" stroke={STK_STROKE} strokeWidth="2" fill="none" opacity="0.55" />
    </svg>
  )
}

export function IconHandshake({ size = 32 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <path d="M11 17l-3-3a3 3 0 010-4l3-3 3 3a3 3 0 010 4l-3 3z" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" />
      <path d="M3 12l4-4M21 12l-4-4M14 14l3 3M10 14l-3 3" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  )
}

export function IconBox({ size = 32 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <path d="M3 7l9-4 9 4M3 7v10l9 4 9-4V7M3 7l9 4M21 7l-9 4M12 11v10" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" />
    </svg>
  )
}

export function IconChart({ size = 32 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <path d="M4 20V10M10 20V4M16 20v-6M22 20H2" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  )
}

export function IconClock({ size = 32 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <circle cx="12" cy="12" r="9" stroke="currentColor" strokeWidth="1.8" />
      <path d="M12 7v5l3 2" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  )
}

export function IconCard({ size = 32 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <rect x="2.5" y="5.5" width="19" height="13" rx="2" stroke="currentColor" strokeWidth="1.8" />
      <path d="M2.5 10h19M6 15h4" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  )
}

export function IconShield({ size = 32 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
      <path d="M12 3l8 3v6c0 5-3.5 8-8 9-4.5-1-8-4-8-9V6l8-3z" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" />
      <path d="M8.5 12l2.5 2.5L16 10" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export function IconCheck({ size = 16 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 16 16" fill="none">
      <path d="M3 8.5L6.5 12L13 5" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export function IconPlus({ size = 16 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 16 16" fill="none">
      <path d="M8 3v10M3 8h10" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" />
    </svg>
  )
}

export function IconArrow({ size = 16 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 16 16" fill="none">
      <path d="M3 8h10M9 4l4 4-4 4" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export function StarIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 100 100">
      <path d="M 50 5 L 56 44 L 95 50 L 56 56 L 50 95 L 44 56 L 5 50 L 44 44 Z" fill="currentColor" />
    </svg>
  )
}

export function StickerApple({ size = 70, tilt = -4 }) {
  return (
    <svg viewBox="0 0 100 100" width={size} height={size} style={{ transform: `rotate(${tilt}deg)` }}>
      <ellipse cx="50" cy="58" rx="34" ry="36" fill="#fff" stroke="#fff" strokeWidth="6" />
      <path d="M 50 24 Q 36 28 30 50 Q 26 76 44 88 Q 50 92 56 88 Q 74 76 70 50 Q 64 28 50 24 Z"
        fill="#D94A4A" stroke="#3D2818" strokeWidth="2.5" />
      <path d="M 50 24 Q 52 18 56 14" stroke="#3D2818" strokeWidth="3" fill="none" strokeLinecap="round" />
      <ellipse cx="60" cy="16" rx="8" ry="4" fill="#4A8B4A" stroke="#3D2818" strokeWidth="2" transform="rotate(28 60 16)" />
      <ellipse cx="40" cy="40" rx="6" ry="3" fill="#fff" opacity="0.4" />
    </svg>
  )
}

export function StickerBag({ size = 80, tilt = 6 }) {
  return (
    <svg viewBox="0 0 100 100" width={size} height={size} style={{ transform: `rotate(${tilt}deg)` }}>
      <rect x="18" y="28" width="64" height="62" rx="6" fill="#fff" stroke="#fff" strokeWidth="6" />
      <rect x="18" y="28" width="64" height="62" rx="6" fill="#4A7C50" stroke="#3D2818" strokeWidth="2.5" />
      <path d="M 32 28 Q 32 10 50 10 Q 68 10 68 28" stroke="#3D2818" strokeWidth="4" fill="none" strokeLinecap="round" />
      <text x="50" y="68" textAnchor="middle" fill="#fff" fontFamily="Unbounded, system-ui" fontWeight="700" fontSize="14">БОКС</text>
    </svg>
  )
}

export function PulseMark({ size = 16, subtle = false }) {
  return (
    <span
      className={`pulse-mark${subtle ? ' pulse-mark--subtle' : ''}`}
      style={{ width: size, height: size }}
      aria-hidden="true"
    >
      <svg viewBox="0 0 24 24" width={size} height={size} fill="none">
        <circle className="pulse-mark__ring" cx="12" cy="13" r="9" fill="currentColor" opacity="0.18" />
        <path className="pulse-mark__leaf"
          d="M5 16c0-6 4-10 14-12-1 9-5 13-11 13-1 0-2-.3-3-1z"
          fill="currentColor" />
        <path d="M5 16C9 12 14 8 18 5"
          stroke="rgba(0,0,0,0.32)" strokeWidth="1" strokeLinecap="round" fill="none" />
        <circle className="pulse-mark__spark" cx="17" cy="6" r="1.4" fill="rgba(255,255,255,0.9)" />
      </svg>
    </span>
  )
}

export function IconPin({ size = 24 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
      <path d="M12 21s-7-7.5-7-12a7 7 0 1 1 14 0c0 4.5-7 12-7 12z" />
      <circle cx="12" cy="9" r="2.5" />
    </svg>
  )
}

export function IconBell({ size = 24 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
      <path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9z" />
      <path d="M10 21a2 2 0 0 0 4 0" />
    </svg>
  )
}

export function IconQR({ size = 24 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round">
      <rect x="3" y="3" width="7" height="7" /><rect x="14" y="3" width="7" height="7" /><rect x="3" y="14" width="7" height="7" />
      <line x1="14" y1="14" x2="14" y2="17" strokeLinecap="round" /><line x1="14" y1="20" x2="14" y2="21" strokeLinecap="round" />
      <line x1="17" y1="14" x2="17" y2="15" strokeLinecap="round" />
      <line x1="20" y1="14" x2="20" y2="17" strokeLinecap="round" />
      <line x1="17" y1="18" x2="21" y2="18" strokeLinecap="round" />
      <line x1="20" y1="21" x2="21" y2="21" strokeLinecap="round" />
    </svg>
  )
}

export function IconEco({ size = 24 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
      <path d="M12 3v18M3 12h18" /><path d="M5 5l14 14M19 5L5 19" strokeWidth="1.2" />
    </svg>
  )
}
