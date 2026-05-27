import { useState, useEffect, useRef } from 'react'
import './landing.css'
import {
  StickerCroissant, StickerCoffee, StickerCoin, StickerBread, StickerSparkle,
  StickerHeart, StickerLeaf, StickerApple, StickerBag,
  PulseMark,
  IconBox, IconClock, IconCard, IconCheck, IconPlus, IconArrow, StarIcon,
  IconPin, IconBell, IconQR, IconEco,
} from './LandingStickers'
import {
  HedgehogHero, HedgehogBox, HedgehogMoney, HedgehogMark,
} from './LandingMascot'
import { PhoneMap, PhoneBoxDetail, PhoneOrdersScreen } from './LandingPhone'

// ─── CountUp ─────────────────────────────────────────────────────────────────
function CountUp({ to, suffix = '', prefix = '', duration = 1400 }) {
  const [val, setVal] = useState(0)
  const ref = useRef(null)
  const startedRef = useRef(false)

  useEffect(() => {
    const obs = new IntersectionObserver((entries) => {
      entries.forEach((e) => {
        if (e.isIntersecting && !startedRef.current) {
          startedRef.current = true
          const start = performance.now()
          const tick = (now) => {
            const t = Math.min(1, (now - start) / duration)
            const eased = 1 - Math.pow(1 - t, 3)
            setVal(Math.round(to * eased))
            if (t < 1) requestAnimationFrame(tick)
          }
          requestAnimationFrame(tick)
        }
      })
    }, { threshold: 0.4 })
    if (ref.current) obs.observe(ref.current)
    return () => obs.disconnect()
  }, [to, duration])

  return <span ref={ref}>{prefix}{val.toLocaleString('ru-RU')}{suffix}</span>
}

// ─── Reveal ───────────────────────────────────────────────────────────────────
function Reveal({ children, delay = 0, className = '' }) {
  const ref = useRef(null)
  const [shown, setShown] = useState(false)

  useEffect(() => {
    const el = ref.current
    if (!el) return
    const r = el.getBoundingClientRect()
    if (r.top < window.innerHeight * 0.9 && r.bottom > 0) {
      const tm = setTimeout(() => setShown(true), Math.min(delay, 50))
      return () => clearTimeout(tm)
    }
    const obs = new IntersectionObserver((es) => {
      es.forEach((e) => {
        if (e.isIntersecting) {
          setTimeout(() => setShown(true), delay)
          obs.disconnect()
        }
      })
    }, { threshold: 0.15 })
    obs.observe(el)
    return () => obs.disconnect()
  }, [delay])

  return (
    <div ref={ref} className={`reveal${shown ? ' in' : ''} ${className}`}>
      {children}
    </div>
  )
}

// ─── SectionMarker ────────────────────────────────────────────────────────────
function SectionMarker({ num, label, dark = false }) {
  return (
    <div className="s-mark" style={dark ? { '--c-ink': 'var(--c-bg)' } : undefined}>
      <span className="s-mark__star">
        <StarIcon />
      </span>
      <span className="s-mark__num" style={dark ? { color: 'var(--c-accent)' } : undefined}>
        №{num}
      </span>
      <span className="s-mark__line" style={dark ? { borderColor: 'var(--c-accent)' } : undefined} />
      <span className="s-mark__lbl" style={dark ? { color: 'var(--c-bg)' } : undefined}>
        {label}
      </span>
    </div>
  )
}

// ─── Parallax hook ────────────────────────────────────────────────────────────
function useParallax(strength = 14) {
  const [off, setOff] = useState({ x: 0, y: 0 })
  useEffect(() => {
    const onMove = (e) => {
      const cx = window.innerWidth / 2
      const cy = window.innerHeight / 2
      setOff({ x: (e.clientX - cx) / cx * strength, y: (e.clientY - cy) / cy * strength })
    }
    window.addEventListener('mousemove', onMove)
    return () => window.removeEventListener('mousemove', onMove)
  }, [strength])
  return off
}

// ─── HeroBackdrop (consumer arc text) ────────────────────────────────────────
function HeroBackdrop() {
  return (
    <div className="hero-stage">
      <div className="hero-stage__blob" />
      <svg className="hero-stage__arc" viewBox="0 0 720 720">
        <defs>
          <path id="hc-hero-arc-top" d="M 110,360 A 250,250 0 0 1 610,360" fill="none" />
        </defs>
        <text>
          <textPath href="#hc-hero-arc-top" startOffset="50%" textAnchor="middle">
            ✦ ЛЮБИМАЯ&nbsp;ЕДА&nbsp;·&nbsp;МЕНЬШЕ&nbsp;ОТХОДОВ&nbsp;·&nbsp;ВЫГОДНЕЕ&nbsp;НА&nbsp;70%&nbsp;✦
          </textPath>
        </text>
      </svg>
      <svg className="hero-stage__doodles" viewBox="0 0 1200 800" preserveAspectRatio="none">
        <path d="M 70 120 q 12 -16 24 0 t 24 0 t 24 0" strokeWidth="2.2" />
        <g transform="translate(180,80)">
          <line x1="-10" y1="0" x2="10" y2="0" strokeWidth="2.2" />
          <line x1="0" y1="-10" x2="0" y2="10" strokeWidth="2.2" />
          <line x1="-7" y1="-7" x2="7" y2="7" strokeWidth="2.2" />
          <line x1="-7" y1="7" x2="7" y2="-7" strokeWidth="2.2" />
        </g>
        <path d="M 540 80 q 30 -30 60 0 t 50 -10" strokeWidth="2.2" strokeDasharray="2 6" />
        <path d="M 1080 140 q -16 -16 -32 -4 q -10 18 14 22 q 22 0 14 -28" strokeWidth="2.2" />
        <g transform="translate(90,520)" strokeWidth="2.4">
          <line x1="-9" y1="0" x2="9" y2="0" />
          <line x1="0" y1="-9" x2="0" y2="9" />
        </g>
        <g transform="translate(280,640)" strokeWidth="0" fill="currentColor">
          <circle cx="0" cy="0" r="3" />
          <circle cx="12" cy="-2" r="3" />
          <circle cx="24" cy="0" r="3" />
        </g>
        <path d="M 800 720 q 18 -18 36 0 t 36 0 t 36 0 t 36 0" strokeWidth="2.2" />
        <g transform="translate(990,640)" strokeWidth="2.4">
          <path d="M 0 -14 L 0 14 M -14 0 L 14 0" />
          <circle cx="0" cy="0" r="3" fill="currentColor" stroke="none" />
        </g>
        <path d="M 30 400 q 30 -20 60 -10 t 50 20" strokeWidth="2.2" />
        <path d="M 130 410 L 142 412 L 138 400" strokeWidth="2.2" />
        <path d="M 1110 340 l 8 -12 l 8 12 l 8 -12 l 8 12" strokeWidth="2.2" />
      </svg>
    </div>
  )
}

// ─── HeroTicker ───────────────────────────────────────────────────────────────
function HeroTicker({ children }) {
  return (
    <span className="ticker">
      <span className="ticker__bar" />
      {children}
      <span className="ticker__bar" />
    </span>
  )
}

// ─── NAV ──────────────────────────────────────────────────────────────────────
function Nav() {
  return (
    <nav className="nav">
      <a href="/" className="nav-brand">
        <span className="brandmark"><HedgehogMark size={36} /></span>
        Бережок
      </a>
      <div className="nav-links">
        <a className="nav-link" href="#how">Как это работает</a>
        <a className="nav-link" href="#cats">Где забрать</a>
        <a className="nav-link" href="#app">Приложение</a>
        <a className="nav-link" href="#faq">FAQ</a>
        <a className="nav-link nav-link--cross" href="/partners">
          <span className="nav-link__dot"><PulseMark size={16} /></span>
          Для бизнеса
        </a>
        <a className="btn btn-primary" href="#download" style={{ marginLeft: 6 }}>
          Скачать
        </a>
      </div>
    </nav>
  )
}

// ─── HERO ─────────────────────────────────────────────────────────────────────
function Hero() {
  const par = useParallax(10)
  return (
    <section className="hero hero-c">
      <HeroBackdrop />
      <div className="wrap hero-grid">
        <div className="hero-left">
          <Reveal>
            <HeroTicker>
              для&nbsp;клиентов <span className="ticker__sep">✷</span> экономия&nbsp;до&nbsp;70%{' '}
              <span className="ticker__sep">✷</span> 500+&nbsp;мест&nbsp;в&nbsp;РФ
            </HeroTicker>
          </Reveal>

          <Reveal delay={80}>
            <h1 className="h-display">
              Любимая еда —{' '}
              <span className="accent-word accent-underline">за&nbsp;треть</span>{' '}
              цены.
            </h1>
          </Reveal>

          <Reveal delay={160}>
            <p className="hero-sub">
              Бережок собирает сюрприз-боксы из остатков любимых пекарен, кофеен
              и ресторанов рядом с вами. То же качество, минус 50–70% от ценника —
              просто заберите вечером.
            </p>
          </Reveal>

          <Reveal delay={240}>
            <div className="hero-ctas">
              <a href="#download" className="btn btn-primary btn-lg">
                Скачать приложение
                <IconArrow size={20} />
              </a>
              <a href="#how" className="btn btn-secondary btn-lg">
                Как это работает
              </a>
            </div>
          </Reveal>

          <Reveal delay={320}>
            <div className="hero-stats">
              <div className="hero-stat">
                <div className="hero-stat-num"><CountUp to={500} suffix="+" /></div>
                <div className="hero-stat-lbl">заведений&nbsp;рядом в&nbsp;7&nbsp;городах</div>
              </div>
              <div className="hero-stat">
                <div className="hero-stat-num"><CountUp to={47} suffix=" т" /></div>
                <div className="hero-stat-lbl">еды&nbsp;спасено пользователями</div>
              </div>
              <div className="hero-stat">
                <div className="hero-stat-num"><CountUp to={1842} suffix="₽" /></div>
                <div className="hero-stat-lbl">средняя&nbsp;экономия в&nbsp;месяц</div>
              </div>
            </div>
          </Reveal>
        </div>

        <div className="hero-right">
          <div
            className="hero-mascot-wrap"
            style={{ '--par-x': `${par.x * -0.6}px`, '--par-y': `${par.y * -0.6}px` }}
          >
            <div className="float-soft" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <PhoneMap />
            </div>
          </div>

          <Reveal delay={420}>
            <div className="hero-c__nudge">
              <span className="hero-c__nudge-dot">МР</span>
              <span>Миша забрал бокс <b>−68%</b><br />2 минуты назад</span>
            </div>
          </Reveal>

          <div className="hero-sticker s-1" style={{ transform: `translate(${par.x * 1.4}px, ${par.y * 1.4}px) rotate(-8deg)` }}>
            <StickerCroissant size={84} tilt={0} />
          </div>
          <div className="hero-sticker s-2" style={{ transform: `translate(${par.x * -1.2}px, ${par.y * -1.2}px) rotate(12deg)` }}>
            <StickerCoffee size={90} tilt={0} />
          </div>
          <div className="hero-sticker s-3" style={{ transform: `translate(${par.x * 1.6}px, ${par.y * 1.6}px) rotate(-12deg)` }}>
            <StickerCoin size={96} tilt={0} />
          </div>
          <div className="hero-sticker s-4" style={{ transform: `translate(${par.x * -1.5}px, ${par.y * -1.5}px) rotate(6deg)` }}>
            <StickerApple size={84} tilt={0} />
          </div>
          <div className="hero-sticker s-5" style={{ transform: `translate(${par.x * 2}px, ${par.y * 2}px) rotate(20deg)` }}>
            <StickerSparkle size={56} tilt={0} />
          </div>
        </div>
      </div>
    </section>
  )
}

// ─── LIVE FEED ────────────────────────────────────────────────────────────────
const FEED_ITEMS = [
  { who: 'Маша', where: '«Корица»', d: '−68%', ago: 'минуту назад' },
  { who: 'Игорь', where: '«Помол»', d: '−55%', ago: '3 минуты назад' },
  { who: 'Аня', where: '«Хлебный двор»', d: '−70%', ago: '5 минут назад' },
  { who: 'Влад', where: '«Шумовка»', d: '−45%', ago: '8 минут назад' },
  { who: 'Лиза', where: '«Café Saturne»', d: '−60%', ago: '11 минут назад' },
  { who: 'Артём', where: '«Чай & Тосты»', d: '−50%', ago: '14 минут назад' },
  { who: 'Катя', where: '«Тёплый плед»', d: '−65%', ago: '16 минут назад' },
  { who: 'Олег', where: '«Маркет 24»', d: '−40%', ago: '19 минут назад' },
]

function LiveFeed() {
  const items = [...FEED_ITEMS, ...FEED_ITEMS]
  return (
    <div className="live-feed">
      <div className="live-feed__mask">
        <div className="live-feed__track">
          {items.map((it, i) => (
            <span key={i} className="live-feed__item">
              <span className="live-feed__dot"><PulseMark size={18} /></span>
              <span className="live-feed__who">{it.who} только что забрал бокс из</span>
              <span className="live-feed__where">{it.where}</span>
              <span className="live-feed__discount">{it.d}</span>
              <span className="live-feed__who">· {it.ago}</span>
              <span className="live-feed__sep">·</span>
            </span>
          ))}
        </div>
      </div>
    </div>
  )
}

// ─── STEPS ────────────────────────────────────────────────────────────────────
function Steps() {
  return (
    <section id="how" className="section steps steps-c">
      <div className="wrap">
        <div className="section-head">
          <SectionMarker num="02" label="Как это работает" />
          <h2 className="h-section">
            Три тапа в приложении —{' '}
            <span className="accent-highlight">ужин в&nbsp;кармане</span>
          </h2>
          <p className="h-sub">
            Не подписки, не доставки, не мутных промокодов. Открыл карту вечером,
            ткнул в заведение по пути домой, забрал бокс — и съел.
          </p>
        </div>

        <div className="steps-grid">
          <Reveal delay={0}>
            <div className="step-card">
              <div className="step-num">1</div>
              <div className="step-illu">
                <HedgehogHero />
                <div style={{ position: 'absolute', top: 12, left: 14 }}>
                  <StickerSparkle size={28} tilt={-12} />
                </div>
              </div>
              <h3>Откройте карту рядом</h3>
              <p>
                Войдите по SMS-коду (без email и пароля), включите геолокацию.
                На карте — все заведения в радиусе 2 км с боксами на сегодня.
              </p>
              <div className="step-meta">
                <IconClock size={14} /> 1 минута на вход
              </div>
            </div>
          </Reveal>

          <Reveal delay={120}>
            <div className="step-card">
              <div className="step-num">2</div>
              <div className="step-illu">
                <HedgehogBox />
                <div style={{ position: 'absolute', bottom: 8, right: 14 }}>
                  <StickerHeart size={26} tilt={14} />
                </div>
              </div>
              <h3>Выберите бокс и оплатите</h3>
              <p>
                В карточке — время выдачи, примерное содержимое, отзывы.
                Оплата в один тап через ЮKassa: карта, СБП, Apple/Google Pay.
              </p>
              <div className="step-meta">
                <IconCard size={14} /> от 120 ₽ за бокс
              </div>
            </div>
          </Reveal>

          <Reveal delay={240}>
            <div className="step-card">
              <div className="step-num">3</div>
              <div className="step-illu">
                <HedgehogMoney />
                <div style={{ position: 'absolute', top: 14, right: 14 }}>
                  <StickerSparkle size={24} tilt={20} />
                </div>
              </div>
              <h3>Заберите в окно выдачи</h3>
              <p>
                Приходите в указанное время, показываете QR-код или 8 символов
                сотруднику — забираете бокс. Никакой кассы и очередей.
              </p>
              <div className="step-meta">
                <IconBox size={14} /> окно 30–60 минут
              </div>
            </div>
          </Reveal>
        </div>
      </div>
    </section>
  )
}

// ─── CATEGORIES ───────────────────────────────────────────────────────────────
const CATS = [
  { t: 'Пекарни',   m: 'Хлеб, выпечка',         c: '120+', sticker: <StickerCroissant size={72} tilt={-6} /> },
  { t: 'Кофейни',   m: 'Кофе, десерты',          c: '180+', sticker: <StickerCoffee size={72} tilt={4} /> },
  { t: 'Рестораны', m: 'Горячее, супы',           c: '95',   sticker: <StickerBread size={72} tilt={-8} /> },
  { t: 'Магазины',  m: 'Готовая еда, продукты',  c: '70',   sticker: <StickerBag size={72} tilt={6} /> },
  { t: 'Отели',     m: 'Завтраки, ланчи',         c: '25',   sticker: <StickerApple size={72} tilt={-4} /> },
  { t: 'Цветочные', m: 'Букеты дня',              c: '18',   sticker: <StickerLeaf size={72} tilt={12} /> },
]

function Categories() {
  return (
    <section id="cats" className="section categories">
      <div className="wrap">
        <div className="section-head">
          <SectionMarker num="03" label="Где забирать" />
          <h2 className="h-section">
            500+ мест{' '}
            <span style={{ fontFamily: 'var(--f-serif)', fontStyle: 'italic', color: 'var(--c-primary)' }}>
              в&nbsp;шаговой
            </span>{' '}
            доступности
          </h2>
          <p className="h-sub">
            Пекарни на углу, кофейни около офиса, ресторанчики у дома. Каждый
            день — новые боксы.
          </p>
        </div>

        <div className="cat-grid">
          {CATS.map((cat, i) => (
            <Reveal key={cat.t} delay={i * 60}>
              <div className="cat-card">
                <span className="cat-card__count">{cat.c}</span>
                <div className="cat-card__pic">{cat.sticker}</div>
                <div>
                  <h3 className="cat-card__t">{cat.t}</h3>
                  <div className="cat-card__meta">{cat.m}</div>
                </div>
              </div>
            </Reveal>
          ))}
        </div>
      </div>
    </section>
  )
}

// ─── BOX REVEAL ───────────────────────────────────────────────────────────────
const BOXES = [
  {
    t: 'Утренний бокс',
    from: 'пекарня «Корица»',
    discount: '−68%',
    time: 'сегодня · 18:00–19:00',
    contents: ['2–3 свежих круассана', 'Сезонные булочки и улитки', 'Печенье или пирожные дня'],
    price: 190, was: 590,
    hero: 'peach',
    icn: '🥐',
  },
  {
    t: 'Вечерний кофе-бокс',
    from: 'сеть «Помол»',
    discount: '−55%',
    time: 'сегодня · 20:00–21:00',
    contents: ['Зерновой кофе 250 г (на выбор)', 'Десерт дня — чизкейк или брауни', 'Маленький снэк-сэндвич'],
    price: 240, was: 530,
    hero: 'moss',
    icn: '☕',
  },
  {
    t: 'Шеф-бокс на двоих',
    from: 'ресторан «Шумовка»',
    discount: '−60%',
    time: 'сегодня · 21:30–22:30',
    contents: ['Горячее блюдо дня', 'Гарнир и салат', 'Свежий хлеб от шефа'],
    price: 590, was: 1480,
    hero: 'sand',
    icn: '🍱',
  },
]

function BoxReveal() {
  return (
    <section className="section boxes">
      <div className="wrap">
        <div className="section-head">
          <SectionMarker num="04" label="Что бывает в боксе" />
          <h2 className="h-section">
            Сюрприз —{' '}
            <span style={{ fontFamily: 'var(--f-serif)', fontStyle: 'italic', color: 'var(--c-primary)' }}>
              но&nbsp;вкусный
            </span>
          </h2>
          <p className="h-sub">
            В боксе — то, что приготовили на сегодня и не успели продать.
            Точный состав узнаете на выдаче. А вот пример того, что бывает внутри:
          </p>
        </div>

        <div className="box-grid">
          {BOXES.map((b, i) => (
            <Reveal key={b.t} delay={i * 100}>
              <div className="box-card">
                <div className={`box-card__hero box-card__hero--${b.hero}`}>
                  <span className="box-card__time">
                    <IconClock size={12} /> {b.time}
                  </span>
                  <span className="box-card__discount">{b.discount}</span>
                  <div style={{ fontSize: 84, position: 'relative', zIndex: 1, lineHeight: 1 }}>{b.icn}</div>
                </div>
                <div className="box-card__body">
                  <h3 className="box-card__t">{b.t}</h3>
                  <div className="box-card__from">из <b>{b.from}</b></div>
                  <ul className="box-card__contents">
                    {b.contents.map((c, j) => <li key={j}>{c}</li>)}
                  </ul>
                  <div className="box-card__price">
                    <b>{b.price}₽</b>
                    <s>{b.was}₽</s>
                  </div>
                </div>
              </div>
            </Reveal>
          ))}
        </div>
      </div>
    </section>
  )
}

// ─── APP SHOWCASE ─────────────────────────────────────────────────────────────
function IconSparkle2({ size = 24 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
      <path d="M12 3v18M3 12h18" />
      <path d="M5 5l14 14M19 5L5 19" strokeWidth="1.2" />
    </svg>
  )
}

function AppShowcase() {
  return (
    <section id="app" className="section app-showcase">
      <div className="wrap">
        <div className="app-top">
          <div>
            <SectionMarker num="05" label="Приложение" dark />
            <h2 className="h-section" style={{ textAlign: 'left', marginTop: 18 }}>
              Карманный фудшеринг —{' '}
              <span style={{ color: 'var(--c-accent)' }}>iOS и Android</span>
            </h2>
            <p className="h-sub" style={{ textAlign: 'left', marginTop: 16, maxWidth: 540 }}>
              Бесплатно, без подписок и скрытых платежей. Открываете карту —
              видите все боксы рядом на сегодня. Заказ за 3 тапа, оплата СБП,
              получение по QR-коду.
            </p>
          </div>

          <Reveal delay={100}>
            <div className="rating-card">
              <div className="rating-card__head">
                <div className="rating-card__big">
                  4.9<span className="denom">/5</span>
                </div>
                <div className="rating-card__head-right">
                  <div className="rating-card__stars">★★★★★</div>
                  <div className="rating-card__meta">
                    более 65 000 оценок<br />в каталогах сторов
                  </div>
                </div>
              </div>

              <div className="rating-card__sources">
                <span className="rating-source">
                  <span className="rating-source__icn">
                    <svg width="12" height="14" viewBox="0 0 24 28" fill="currentColor">
                      <path d="M19.6 14.1c0-3.7 3-5.4 3.1-5.5-1.7-2.5-4.4-2.8-5.3-2.9-2.3-.2-4.4 1.3-5.6 1.3-1.2 0-2.9-1.3-4.8-1.2-2.5 0-4.7 1.4-6 3.6-2.6 4.5-.7 11 1.8 14.6 1.2 1.8 2.6 3.7 4.6 3.6 1.8-.1 2.5-1.2 4.7-1.2s2.8 1.2 4.8 1.2c2 0 3.2-1.8 4.4-3.5 1.4-2 2-3.9 2-4-.1 0-3.8-1.4-3.7-5.6z" />
                      <path d="M16 4c1-1.2 1.7-2.9 1.5-4.5-1.5.1-3.2 1-4.3 2.2-1 1-1.8 2.7-1.6 4.4 1.6.1 3.4-.9 4.4-2.1z" />
                    </svg>
                  </span>
                  App Store <b>4.9</b>
                </span>
                <span className="rating-source">
                  <span className="rating-source__icn">
                    <svg width="12" height="14" viewBox="0 0 24 26" fill="currentColor">
                      <path d="M3.3.5L14.7 11.8l-2.1 2.1L3.3 25c-.5-.1-.8-.5-.8-1.1V1.6c0-.5.3-1 .8-1.1z" />
                      <path d="M16.7 16.1l2.9-2.9c.9-.5.9-1.9 0-2.4l-2.9-2.9-3.1 3.1 3.1 3.1z" />
                    </svg>
                  </span>
                  Google Play <b>4.8</b>
                </span>
                <span className="rating-source">
                  <span className="rating-source__icn">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 0L8 4l4 4 4-4-4-4zM4 8l-4 4 4 4 4-4-4-4zm16 0l-4 4 4 4 4-4-4-4zM12 16l-4 4 4 4 4-4-4-4z" />
                    </svg>
                  </span>
                  RuStore <b>4.9</b>
                </span>
              </div>

              <div className="rating-card__divider" />
              <div className="rating-quote">
                <div className="rating-quote__avatar">КЛ</div>
                <div className="rating-quote__body">
                  <p className="rating-quote__t">
                    «Скачала ради эксперимента — за месяц забрала 12 боксов
                    и сэкономила почти 6000₽. И ни одного разочарования.»
                  </p>
                  <div className="rating-quote__who">
                    <b>Ксения Л.</b> · App Store · 14 мая
                  </div>
                </div>
              </div>
            </div>
          </Reveal>
        </div>

        <div className="app-facts">
          <div className="app-fact">
            <div className="app-fact__n">3 тапа</div>
            <div className="app-fact__l">от&nbsp;карты до&nbsp;брони</div>
          </div>
          <div className="app-fact">
            <div className="app-fact__n">~42 МБ</div>
            <div className="app-fact__l">размер установки</div>
          </div>
          <div className="app-fact">
            <div className="app-fact__n">7 городов</div>
            <div className="app-fact__l">и+&nbsp;каждый месяц новые</div>
          </div>
          <div className="app-fact">
            <div className="app-fact__n">0 ₽</div>
            <div className="app-fact__l">без&nbsp;подписки и&nbsp;скрытых плат</div>
          </div>
        </div>

        <div className="fan-wrap">
          <div className="fan-phones">
            <PhoneOrdersScreen />
            <PhoneMap />
            <PhoneBoxDetail />
          </div>

          <div className="feature-list">
            <Reveal>
              <div className="feature">
                <div className="feature__icn"><IconPin size={26} /></div>
                <div>
                  <h3 className="feature__t">Карта вокруг вас</h3>
                  <p className="feature__d">
                    Все активные боксы на Яндекс.Картах. Фильтры по категории, цене,
                    расстоянию и времени выдачи.
                  </p>
                </div>
              </div>
            </Reveal>
            <Reveal delay={100}>
              <div className="feature">
                <div className="feature__icn"><IconBell size={26} /></div>
                <div>
                  <h3 className="feature__t">Смарт-уведомления</h3>
                  <p className="feature__d">
                    Тихий пуш, когда любимая пекарня выкладывает бокс. Никакого спама —
                    только то, что вы пометили звёздочкой.
                  </p>
                </div>
              </div>
            </Reveal>
            <Reveal delay={200}>
              <div className="feature">
                <div className="feature__icn"><IconQR size={24} /></div>
                <div>
                  <h3 className="feature__t">QR-код на выдаче</h3>
                  <p className="feature__d">
                    Покажите код сотруднику — забрали бокс. Без кассы и очередей.
                    Если опоздали — деньги вернутся.
                  </p>
                </div>
              </div>
            </Reveal>
            <Reveal delay={300}>
              <div className="feature">
                <div className="feature__icn"><IconSparkle2 size={24} /></div>
                <div>
                  <h3 className="feature__t">Личный эко-счётчик</h3>
                  <p className="feature__d">
                    Сколько еды и денег вы спасли, на сколько обогнали друзей. Делитесь
                    скриншотом, копите бейджи.
                  </p>
                </div>
              </div>
            </Reveal>
          </div>
        </div>
      </div>
    </section>
  )
}

// ─── ECO IMPACT ───────────────────────────────────────────────────────────────
function EcoImpact() {
  return (
    <section className="section eco">
      <div className="wrap">
        <div className="eco-grid">
          <Reveal>
            <div>
              <SectionMarker num="06" label="Эко-эффект" />
              <h2 className="h-section" style={{ textAlign: 'left', marginTop: 18 }}>
                Каждый бокс —{' '}
                <span style={{ fontFamily: 'var(--f-serif)', fontStyle: 'italic', color: 'var(--c-primary)' }}>
                  минус&nbsp;1.4&nbsp;кг
                </span>{' '}
                CO₂
              </h2>
              <p className="h-sub" style={{ textAlign: 'left', marginTop: 16 }}>
                Пищевые отходы — это 8% мировых выбросов парниковых газов.
                Один бокс — это около килограмма еды, которая не поедет на свалку.
                В приложении вы видите свой личный эко-счёт и сравниваете его с друзьями.
              </p>
              <div style={{ display: 'flex', gap: 12, marginTop: 28, alignItems: 'center' }}>
                <StickerLeaf size={48} tilt={-10} />
                <StickerSparkle size={32} tilt={20} />
                <span style={{ fontFamily: 'var(--f-mono)', fontSize: 12, color: 'var(--c-ink-soft)', textTransform: 'uppercase', letterSpacing: '0.12em' }}>
                  пример экрана →
                </span>
              </div>
            </div>
          </Reveal>

          <Reveal delay={120}>
            <div className="eco-card">
              <div className="eco-card__head">
                <div>
                  <div className="eco-card__name">Ваш эко-счёт</div>
                  <div style={{ fontFamily: 'var(--f-display)', fontWeight: 700, fontSize: 18, marginTop: 4 }}>
                    Маша, привет 👋
                  </div>
                </div>
                <span className="eco-card__tier">
                  <span className="eco-card__tier-name">Хранитель</span>
                  <span className="eco-card__tier-pips" aria-label="Уровень 4 из 6">
                    {[0, 1, 2, 3, 4, 5].map((i) => (
                      <svg
                        key={i}
                        className={`eco-card__tier-pip${i < 4 ? ' is-on' : ''}`}
                        width="11" height="13" viewBox="0 0 14 16" fill="currentColor"
                      >
                        <path d="M3 11c0-4 3-7 9-8-1 6-3 9-7 9-1 0-1.5-.3-2-1z" />
                      </svg>
                    ))}
                  </span>
                </span>
              </div>

              <div className="eco-card__big">
                42.3
                <span className="small">кг спасено</span>
              </div>

              <p className="eco-card__sub">
                Это примерно как 28 обедов, которые могли бы поехать в мусор.
                Вы экономите 1&thinsp;820&nbsp;₽ в месяц и обходите 84% пользователей.
              </p>

              <div className="eco-stats">
                <div className="eco-stats__cell">
                  <div className="eco-stats__n"><CountUp to={47} /></div>
                  <div className="eco-stats__l">боксов забрано</div>
                </div>
                <div className="eco-stats__cell">
                  <div className="eco-stats__n"><CountUp to={59} suffix=" кг" /></div>
                  <div className="eco-stats__l">CO₂ не выпущено</div>
                </div>
                <div className="eco-stats__cell">
                  <div className="eco-stats__n"><CountUp to={12420} suffix=" ₽" /></div>
                  <div className="eco-stats__l">экономия за год</div>
                </div>
              </div>
            </div>
          </Reveal>
        </div>
      </div>
    </section>
  )
}

// ─── REVIEWS ──────────────────────────────────────────────────────────────────
const C_REVIEWS = [
  {
    initials: 'МО',
    name: 'Маша О.',
    role: 'студент, СПб',
    quote: 'Беру бокс из пекарни рядом с общагой — 180₽ вместо 600. Хватает на завтрак и перекус на парах. Скиньте такое в каждый город страны.',
    metric: '−68% на завтрак',
    stars: 5,
  },
  {
    initials: 'ИК',
    name: 'Иван К.',
    role: 'разработчик, Москва',
    quote: 'Беру вечером после работы — нашёл рамен-боксы у одного места, теперь два раза в неделю минимум. И едой не разбрасываются, и я в плюсе.',
    metric: '2 раза/нед',
    stars: 5,
  },
  {
    initials: 'АТ',
    name: 'Алиса Т.',
    role: 'копирайтер, Казань',
    quote: 'Главное — это что я перестала виновато выбрасывать недоеденные круассаны из соседней пекарни. Теперь те же круассаны живут у меня. Это просто справедливо.',
    metric: 'уже 6 месяцев',
    stars: 5,
  },
]

function Reviews() {
  return (
    <section className="section reviews">
      <div className="wrap">
        <div className="section-head">
          <SectionMarker num="07" label="Что говорят" />
          <h2 className="h-section">
            «Попробую один раз» —{' '}
            <span style={{ fontFamily: 'var(--f-serif)', fontStyle: 'italic', color: 'var(--c-primary)' }}>
              через&nbsp;месяц
            </span>{' '}
            не выходят из приложения
          </h2>
        </div>

        <div className="reviews-grid">
          {C_REVIEWS.map((r, i) => (
            <Reveal key={r.name} delay={i * 100}>
              <div className="review-card">
                <span className="review-metric">{r.metric}</span>
                <div className="review-stars">{'★'.repeat(r.stars)}</div>
                <p className="review-quote">«{r.quote}»</p>
                <div className="review-footer">
                  <div className="review-avatar">{r.initials}</div>
                  <div className="review-meta">
                    <div className="review-name">{r.name}</div>
                    <div className="review-role">{r.role}</div>
                  </div>
                </div>
              </div>
            </Reveal>
          ))}
        </div>
      </div>
    </section>
  )
}

// ─── FAQ ──────────────────────────────────────────────────────────────────────
const C_FAQS = [
  {
    q: 'А что конкретно будет в боксе? Я не люблю сюрпризы.',
    a: 'В карточке заведения всегда написано примерное содержимое — «2–3 круассана и булочки», «горячее блюдо дня и гарнир». Точный состав вы видите при получении. Если что-то категорически не подходит — можно вернуть в течение 15 минут.',
  },
  {
    q: 'Еда же остатки — она ещё свежая?',
    a: 'Это не остатки в смысле «вчерашнее». Это продукция этого дня, которая не успела продаться к вечеру. Партнёры — те же пекарни и кафе, в которые вы ходите днём, только со скидкой к закрытию.',
  },
  {
    q: 'Можно ли отменить заказ или вернуть деньги?',
    a: 'Если партнёр не подтвердил заказ в течение 30 минут — деньги возвращаются автоматически. Если вы не успели прийти — заказ сгорает (но в течение 15 минут после сканирования можно открыть спор, если что-то пошло не так).',
  },
  {
    q: 'Сколько стоит приложение и есть ли подписка?',
    a: 'Приложение бесплатное, без подписки. Вы платите только за бокс, который заказали. Никакой абонентки, ограничений на количество заказов или премиум-тарифа нет.',
  },
  {
    q: 'В моём городе есть заведения?',
    a: 'Прямо сейчас работаем в 7 городах: Москва, СПб, Казань, Екатеринбург, Новосибирск, Краснодар, Нижний Новгород. Города добавляем по мере подключения партнёров — в приложении можно подписаться на уведомление о старте.',
  },
  {
    q: 'А как с оплатой? Безопасно?',
    a: 'Оплата через ЮKassa: карты Мир/Visa/Mastercard, СБП, Apple/Google Pay. Бережок не хранит данные карт — всё проходит через защищённую платёжную систему.',
  },
  {
    q: 'Что если я опоздал на окно выдачи?',
    a: 'Окно — обычно 30–60 минут (зависит от заведения). За час до его окончания приходит напоминание. Если опоздали — увы, бокс сгорает: партнёр уже его собрал и держал именно для вас.',
  },
]

function FAQItem({ q, a, open, onClick }) {
  return (
    <div className={`faq-item${open ? ' open' : ''}`} onClick={onClick}>
      <div className="faq-q">
        <span>{q}</span>
        <span className="faq-icn"><IconPlus size={16} /></span>
      </div>
      <div className="faq-a">{a}</div>
    </div>
  )
}

function FAQ() {
  const [open, setOpen] = useState(0)
  return (
    <section id="faq" className="section faq">
      <div className="wrap">
        <div className="faq-grid">
          <div className="faq-aside">
            <div style={{ marginBottom: 18 }}>
              <SectionMarker num="08" label="FAQ" />
            </div>
            <h2 className="h-section" style={{ textAlign: 'left' }}>
              Что обычно спрашивают перед первым заказом
            </h2>
            <p className="h-sub" style={{ textAlign: 'left' }}>
              Если ответа не нашли — напишите в чат поддержки в приложении.
              Ответим живым человеком, обычно за 5–10 минут.
            </p>
            <div style={{ marginTop: 28, display: 'flex', gap: 12 }}>
              <StickerLeaf size={48} tilt={-10} />
              <StickerHeart size={40} tilt={12} />
              <StickerSparkle size={32} tilt={20} />
            </div>
          </div>
          <div className="faq-list">
            {C_FAQS.map((f, i) => (
              <FAQItem
                key={i}
                q={f.q}
                a={f.a}
                open={open === i}
                onClick={() => setOpen(open === i ? -1 : i)}
              />
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}

// ─── DOWNLOAD CARD ────────────────────────────────────────────────────────────
function DownloadCard() {
  const [phone, setPhone] = useState('')
  return (
    <div className="download-card">
      <h3 className="download-card__t">Скачать Бережок</h3>
      <div className="download-card__row">
        <div className="download-card__qr">
          <svg viewBox="0 0 50 50" xmlns="http://www.w3.org/2000/svg">
            <rect width="50" height="50" fill="var(--c-primary-deep)" />
            {[
              [3,3,10,10],[37,3,10,10],[3,37,10,10],
              [5,5,6,6],[39,5,6,6],[5,39,6,6],
              [15,3,2,2],[19,3,4,2],[25,3,2,4],[29,3,2,2],[33,3,4,2],
              [3,15,2,2],[3,19,2,4],[3,25,4,2],[3,31,2,2],[3,35,2,2],
              [15,15,4,4],[21,15,2,2],[25,15,2,4],[29,15,4,2],[35,15,2,4],
              [15,21,2,4],[19,23,4,2],[25,21,4,4],[31,21,2,2],[35,23,4,2],
              [15,29,4,2],[21,29,2,4],[27,29,4,2],[33,29,2,4],[39,29,2,2],
              [15,35,2,4],[19,37,4,2],[25,35,2,2],[29,35,4,4],[35,35,2,4],[39,39,4,4],
              [23,41,2,2],[27,41,2,4],[33,43,2,2],
            ].map(([x, y, w, h], i) => (
              <rect key={i} x={x} y={y} width={w} height={h} fill="var(--c-bg)" />
            ))}
          </svg>
        </div>
        <div className="download-card__qrhint">
          <b>Наведите камеру телефона</b> на код —<br />
          откроется страница в&nbsp;вашем сторе.
        </div>
      </div>

      <div className="store-badges">
        <a className="store-badge" href="#">
          <svg width="22" height="26" viewBox="0 0 24 28" fill="currentColor" aria-hidden="true">
            <path d="M19.6 14.1c0-3.7 3-5.4 3.1-5.5-1.7-2.5-4.4-2.8-5.3-2.9-2.3-.2-4.4 1.3-5.6 1.3-1.2 0-2.9-1.3-4.8-1.2-2.5 0-4.7 1.4-6 3.6-2.6 4.5-.7 11 1.8 14.6 1.2 1.8 2.6 3.7 4.6 3.6 1.8-.1 2.5-1.2 4.7-1.2s2.8 1.2 4.8 1.2c2 0 3.2-1.8 4.4-3.5 1.4-2 2-3.9 2-4-.1 0-3.8-1.4-3.7-5.6z" />
            <path d="M16 4c1-1.2 1.7-2.9 1.5-4.5-1.5.1-3.2 1-4.3 2.2-1 1-1.8 2.7-1.6 4.4 1.6.1 3.4-.9 4.4-2.1z" />
          </svg>
          <span>
            <div className="store-badge__small">Скачать в</div>
            <div className="store-badge__big">App Store</div>
          </span>
        </a>
        <a className="store-badge" href="#">
          <svg width="22" height="24" viewBox="0 0 24 26" fill="currentColor" aria-hidden="true">
            <path d="M2 1.5L13.6 13 2 24.5c-.3-.4-.5-.9-.5-1.4V2.9C1.5 2.4 1.7 1.9 2 1.5z" />
            <path d="M16.7 16.1l2.9-2.9c.9-.5.9-1.9 0-2.4l-2.9-2.9-3.1 3.1 3.1 3.1z" />
            <path d="M3.3.5L14.7 11.8l-2.1 2.1L3.3 25c-.5-.1-.8-.5-.8-1.1V1.6c0-.5.3-1 .8-1.1z" />
            <path d="M3.3 25L14.7 13.7 16.7 16.1c-.4.2-12.5 8.7-12.5 8.7-.3.2-.6.2-.9.2z" />
          </svg>
          <span>
            <div className="store-badge__small">Доступно в</div>
            <div className="store-badge__big">Google Play</div>
          </span>
        </a>
        <a className="store-badge store-badge--ru" href="#">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path d="M12 0L8 4l4 4 4-4-4-4zM4 8l-4 4 4 4 4-4-4-4zm16 0l-4 4 4 4 4-4-4-4zM12 16l-4 4 4 4 4-4-4-4z" />
          </svg>
          <span>
            <div className="store-badge__small">Также в</div>
            <div className="store-badge__big">RuStore</div>
          </span>
        </a>
      </div>

      <div className="cta-sms">
        <input
          type="tel"
          placeholder="+7 (___) ___-__-__"
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
        />
        <button type="button">Прислать ссылку</button>
      </div>
    </div>
  )
}

// ─── FINAL SECTION ────────────────────────────────────────────────────────────
function Final() {
  return (
    <section id="download" className="final final-c">
      <div className="wrap final-grid">
        <Reveal>
          <div>
            <div style={{ marginBottom: 22 }}>
              <SectionMarker num="09" label="Скачать" dark />
            </div>
            <h2 className="h-section final-h">
              Поужинать{' '}
              <span style={{ color: 'var(--c-accent)' }}>выгоднее</span>{' '}
              можно прямо сегодня.
            </h2>
            <p className="final-sub">
              Скачайте приложение, войдите по SMS — и в радиусе 2 км вокруг вас
              откроются десятки заведений с боксами на сегодня вечером.
            </p>
            <div className="final-checks">
              <div className="final-check">
                <span className="final-check-icn"><IconCheck size={14} /></span>
                <span>Без подписок и абонентки — платите только за бокс</span>
              </div>
              <div className="final-check">
                <span className="final-check-icn"><IconCheck size={14} /></span>
                <span>Регистрация по SMS — без email, паролей и анкет</span>
              </div>
              <div className="final-check">
                <span className="final-check-icn"><IconCheck size={14} /></span>
                <span>Оплата картой, СБП, Apple/Google Pay</span>
              </div>
              <div className="final-check">
                <span className="final-check-icn"><IconCheck size={14} /></span>
                <span>Деньги вернутся, если что-то пошло не так</span>
              </div>
            </div>
          </div>
        </Reveal>

        <Reveal delay={150}>
          <DownloadCard />
        </Reveal>
      </div>
    </section>
  )
}

// ─── FOOTER ───────────────────────────────────────────────────────────────────
function Footer() {
  return (
    <footer className="footer">
      <div className="footer-brand">
        <HedgehogMark size={28} />
        Бережок
      </div>
      <div>© 2026 Бережок · Меньше отходов — больше добра.</div>
      <div style={{ display: 'flex', gap: 18 }}>
        <a style={{ color: 'inherit' }} href="/partners">Для бизнеса</a>
        <a style={{ color: 'inherit' }} href="#">Условия</a>
        <a style={{ color: 'inherit' }} href="#">Конфиденциальность</a>
        <a style={{ color: 'inherit' }} href="#">help@berezhok.ru</a>
      </div>
    </footer>
  )
}

// ─── PAGE ─────────────────────────────────────────────────────────────────────
export default function HomePage() {
  return (
    <div className="lp page--consumer">
      <Nav />
      <Hero />
      <LiveFeed />
      <Steps />
      <Categories />
      <BoxReveal />
      <AppShowcase />
      <EcoImpact />
      <Reviews />
      <FAQ />
      <Final />
      <Footer />
    </div>
  )
}
