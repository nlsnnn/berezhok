import { useState, useEffect, useRef } from 'react'
import { IconArrow } from './LandingStickers'

const CATEGORIES = [
  { id: 'bakery',  label: 'Пекарня',  price: 250, boxes: 8,  emoji: '🥐' },
  { id: 'cafe',    label: 'Кафе',     price: 350, boxes: 5,  emoji: '☕' },
  { id: 'rest',    label: 'Ресторан', price: 550, boxes: 4,  emoji: '🍝' },
  { id: 'shop',    label: 'Магазин',  price: 400, boxes: 12, emoji: '🛒' },
  { id: 'hotel',   label: 'Отель',    price: 700, boxes: 6,  emoji: '🏨' },
  { id: 'flowers', label: 'Цветы',    price: 600, boxes: 3,  emoji: '💐' },
]

function useAnimatedNumber(target, duration = 600) {
  const [val, setVal] = useState(target)
  const prevRef = useRef(target)
  const rafRef = useRef(null)

  useEffect(() => {
    const start = performance.now()
    const from = prevRef.current
    const to = target
    const tick = (now) => {
      const t = Math.min(1, (now - start) / duration)
      const eased = 1 - Math.pow(1 - t, 3)
      const v = from + (to - from) * eased
      setVal(v)
      if (t < 1) rafRef.current = requestAnimationFrame(tick)
      else prevRef.current = to
    }
    rafRef.current = requestAnimationFrame(tick)
    return () => cancelAnimationFrame(rafRef.current)
  }, [target, duration])

  return val
}

function fmtRub(n) {
  return Math.round(n).toLocaleString('ru-RU') + ' ₽'
}

export default function LandingCalculator({ SectionMarker }) {
  const [catId, setCatId] = useState('bakery')
  const [price, setPrice] = useState(250)
  const [boxes, setBoxes] = useState(8)
  const [days, setDays] = useState(26)

  const onPick = (id) => {
    const c = CATEGORIES.find((c) => c.id === id)
    setCatId(id)
    if (c) { setPrice(c.price); setBoxes(c.boxes) }
  }

  const daily = price * boxes
  const monthlyGross = daily * days
  const promoComm = 0.10
  const monthlyNet = monthlyGross * (1 - promoComm)
  const yearlyNet = monthlyNet * 4 + (monthlyGross * (1 - 0.20)) * 8
  const savedWaste = boxes * days * 0.45

  const aMonth = useAnimatedNumber(monthlyNet)
  const aDay = useAnimatedNumber(daily * (1 - promoComm))
  const aYear = useAnimatedNumber(yearlyNet)
  const aWaste = useAnimatedNumber(savedWaste)

  const pricePct = ((price - 100) / (1500 - 100)) * 100
  const boxesPct = ((boxes - 1) / (25 - 1)) * 100
  const daysPct = ((days - 10) / (30 - 10)) * 100

  return (
    <div className="calc-wrap">
      <div className="calc-controls">
        <div>
          {SectionMarker}
          <h3 style={{
            fontFamily: 'var(--f-display)',
            fontWeight: 700,
            fontSize: '32px',
            lineHeight: 1.05,
            margin: '18px 0 12px',
            letterSpacing: '-0.02em',
          }}>
            Сколько вы заработаете на остатках?
          </h3>
          <p style={{ margin: 0, color: 'var(--c-ink-soft)', fontSize: '15px' }}>
            Выберите тип заведения и подвигайте ползунки — посчитаем за 2 секунды
          </p>
        </div>

        <div className="calc-categories">
          {CATEGORIES.map((c) => (
            <button
              key={c.id}
              className={`calc-cat${catId === c.id ? ' active' : ''}`}
              onClick={() => onPick(c.id)}
            >
              <span style={{ fontSize: 15 }}>{c.emoji}</span>
              {c.label}
            </button>
          ))}
        </div>

        <div className="slider-block">
          <div className="slider-head">
            <label htmlFor="s-price">Цена сюрприз-бокса</label>
            <span className="slider-val">{price} ₽</span>
          </div>
          <input
            id="s-price"
            type="range" min="100" max="1500" step="50"
            className="slider"
            style={{ '--pct': `${pricePct}%` }}
            value={price}
            onChange={(e) => setPrice(+e.target.value)}
          />
          <div className="slider-range"><span>100 ₽</span><span>1500 ₽</span></div>
        </div>

        <div className="slider-block">
          <div className="slider-head">
            <label htmlFor="s-boxes">Боксов в день</label>
            <span className="slider-val">{boxes} шт</span>
          </div>
          <input
            id="s-boxes"
            type="range" min="1" max="25" step="1"
            className="slider"
            style={{ '--pct': `${boxesPct}%` }}
            value={boxes}
            onChange={(e) => setBoxes(+e.target.value)}
          />
          <div className="slider-range"><span>1</span><span>25</span></div>
        </div>

        <div className="slider-block">
          <div className="slider-head">
            <label htmlFor="s-days">Дней работы в месяц</label>
            <span className="slider-val">{days}</span>
          </div>
          <input
            id="s-days"
            type="range" min="10" max="30" step="1"
            className="slider"
            style={{ '--pct': `${daysPct}%` }}
            value={days}
            onChange={(e) => setDays(+e.target.value)}
          />
          <div className="slider-range"><span>10</span><span>30</span></div>
        </div>
      </div>

      <div className="calc-output">
        <h3>Ваш доход с Бережка</h3>

        <div className="calc-headline">
          {fmtRub(aMonth)}<span className="small">/мес</span>
        </div>

        <div>
          <div className="calc-row">
            <span className="lbl">В день после комиссии</span>
            <span className="v">{fmtRub(aDay)}</span>
          </div>
          <div className="calc-row">
            <span className="lbl">За год (с учётом промо)</span>
            <span className="v">{fmtRub(aYear)}</span>
          </div>
          <div className="calc-row">
            <span className="lbl">Спасли еды от мусорки</span>
            <span className="v">{aWaste.toFixed(0)} кг/мес</span>
          </div>
        </div>

        <div className="calc-tip">
          ↳ Комиссия 10% в первые 3 месяца, далее 20%.
          Это деньги, которые иначе вы бы просто выбросили.
        </div>

        <a href="#apply" className="btn btn-accent btn-lg" style={{ alignSelf: 'flex-start', marginTop: 6 }}>
          Стать партнёром
          <IconArrow size={18} />
        </a>
      </div>
    </div>
  )
}
