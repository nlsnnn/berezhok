import { useState, useEffect, useRef } from 'react'
import './landing.css'
import {
  StickerCroissant, StickerCoffee, StickerCoin, StickerBread, StickerSparkle,
  StickerHeart, StickerLeaf,
  IconHandshake, IconBox, IconChart, IconClock, IconCard, IconShield,
  IconCheck, IconPlus, IconArrow, StarIcon,
} from './LandingStickers'
import {
  HedgehogHero, HedgehogApply, HedgehogBox, HedgehogMoney, HedgehogMark,
} from './LandingMascot'
import LandingCalculator from './LandingCalculator'
import LandingFinalForm from './LandingFinalForm'

// ─── CountUp (count-up on intersection) ─────────────────────────────────────
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

// ─── Reveal (fade-in on scroll) ──────────────────────────────────────────────
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

// ─── SectionMarker ───────────────────────────────────────────────────────────
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

// ─── HeroBackdrop ────────────────────────────────────────────────────────────
function HeroBackdrop() {
  return (
    <div className="hero-stage">
      <div className="hero-stage__blob" />
      <svg className="hero-stage__arc" viewBox="0 0 720 720">
        <defs>
          <path id="lp-hero-arc-top" d="M 110,360 A 250,250 0 0 1 610,360" fill="none" />
        </defs>
        <text>
          <textPath href="#lp-hero-arc-top" startOffset="50%" textAnchor="middle">
            ✦ ХРАНИТЕЛЬ&nbsp;ОСТАТКОВ&nbsp;·&nbsp;ЗАЩИТНИК&nbsp;ЕДЫ&nbsp;·&nbsp;EST.&nbsp;2026&nbsp;✦
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
        <path d="M 700 30 q 12 -22 24 -8 q 8 14 -10 18 q -16 -2 -2 -14" strokeWidth="2.2" />
      </svg>
    </div>
  )
}

// ─── Nav ─────────────────────────────────────────────────────────────────────
function Nav() {
  return (
    <nav className="nav">
      <a href="#" className="nav-brand">
        <HedgehogMark size={36} />
        Бережок
      </a>
      <div className="nav-links">
        <a className="nav-link" href="#how">Как это работает</a>
        <a className="nav-link" href="#calc">Калькулятор</a>
        <a className="nav-link" href="#reviews">Партнёры</a>
        <a className="nav-link" href="#faq">FAQ</a>
        <a className="btn btn-primary" href="#apply" style={{ marginLeft: 6 }}>
          Стать партнёром
        </a>
      </div>
    </nav>
  )
}

// ─── Hero ─────────────────────────────────────────────────────────────────────
function Hero() {
  const [par, setPar] = useState({ x: 0, y: 0 })

  useEffect(() => {
    const onMove = (e) => {
      const cx = window.innerWidth / 2
      const cy = window.innerHeight / 2
      const nx = (e.clientX - cx) / cx
      const ny = (e.clientY - cy) / cy
      setPar({ x: nx * 14, y: ny * 14 })
    }
    window.addEventListener('mousemove', onMove)
    return () => window.removeEventListener('mousemove', onMove)
  }, [])

  return (
    <section className="hero">
      <HeroBackdrop />
      <div className="wrap hero-grid">
        <div className="hero-left">
          <Reveal>
            <span className="ticker">
              <span className="ticker__bar" />
              для&nbsp;партнёров <span className="ticker__sep">✷</span> 500+&nbsp;заведений <span className="ticker__sep">✷</span> работаем&nbsp;в&nbsp;РФ
              <span className="ticker__bar" />
            </span>
          </Reveal>

          <Reveal delay={80}>
            <h1 className="h-display">
              Не выбрасывайте.{' '}
              <span className="accent-word accent-underline">Продавайте</span>{' '}
              остатки.
            </h1>
          </Reveal>

          <Reveal delay={160}>
            <p className="hero-sub">
              Бережок превращает непроданную еду в&nbsp;дополнительную выручку.
              Загрузите остатки в&nbsp;сюрприз-бокс — клиент заберёт в&nbsp;конце дня
              и оставит у вас деньги, а не в мусорке.
            </p>
          </Reveal>

          <Reveal delay={240}>
            <div className="hero-ctas">
              <a href="#apply" className="btn btn-primary btn-lg">
                Стать партнёром <IconArrow size={20} />
              </a>
              <a href="#calc" className="btn btn-secondary btn-lg">
                Посчитать доход
              </a>
            </div>
          </Reveal>

          <Reveal delay={320}>
            <div className="hero-stats">
              <div className="hero-stat">
                <div className="hero-stat-num"><CountUp to={500} suffix="+" /></div>
                <div className="hero-stat-lbl">партнёров&nbsp;уже&nbsp;с&nbsp;нами</div>
              </div>
              <div className="hero-stat">
                <div className="hero-stat-num"><CountUp to={47} suffix=" т" /></div>
                <div className="hero-stat-lbl">еды&nbsp;спасено от&nbsp;мусорки</div>
              </div>
              <div className="hero-stat">
                <div className="hero-stat-num"><CountUp to={32} suffix="%" /></div>
                <div className="hero-stat-lbl">средний рост выручки&nbsp;к&nbsp;закрытию</div>
              </div>
            </div>
          </Reveal>
        </div>

        <div className="hero-right">
          <div
            className="hero-mascot-wrap"
            style={{ '--par-x': `${par.x * -0.6}px`, '--par-y': `${par.y * -0.6}px` }}
          >
            <div className="float-soft">
              <HedgehogHero pose="hi" />
            </div>
          </div>

          <div className="hero-sticker s-1" style={{
            '--tilt': '-8deg',
            transform: `translate(${par.x * 1.5}px, ${par.y * 1.5}px) rotate(-8deg)`,
          }}>
            <StickerCroissant size={92} tilt={0} />
          </div>
          <div className="hero-sticker s-2" style={{
            '--tilt': '12deg',
            transform: `translate(${par.x * -1.2}px, ${par.y * -1.2}px) rotate(12deg)`,
          }}>
            <StickerCoffee size={88} tilt={0} />
          </div>
          <div className="hero-sticker s-3" style={{
            '--tilt': '-12deg',
            transform: `translate(${par.x * 1.8}px, ${par.y * 1.8}px) rotate(-12deg)`,
          }}>
            <StickerCoin size={104} tilt={0} />
          </div>
          <div className="hero-sticker s-4" style={{
            '--tilt': '6deg',
            transform: `translate(${par.x * -1.6}px, ${par.y * -1.6}px) rotate(6deg)`,
          }}>
            <StickerBread size={92} tilt={0} />
          </div>
          <div className="hero-sticker s-5" style={{
            '--tilt': '20deg',
            transform: `translate(${par.x * 2}px, ${par.y * 2}px) rotate(20deg)`,
          }}>
            <StickerSparkle size={56} tilt={0} />
          </div>
        </div>
      </div>
    </section>
  )
}

// ─── Steps ───────────────────────────────────────────────────────────────────
function Steps() {
  return (
    <section id="how" className="section steps">
      <div className="wrap">
        <div className="section-head">
          <SectionMarker num="02" label="Как это работает" />
          <h2 className="h-section">
            Три шага — и заработаете{' '}
            <span className="accent-highlight">уже сегодня</span>
          </h2>
          <p className="h-sub">
            Никакого софта на кассе, договоров на сто страниц и менеджеров с КПД ноль.
            Зарегистрировались за 5 минут — приняли первый заказ к ужину.
          </p>
        </div>

        <div className="steps-grid">
          <Reveal delay={0}>
            <div className="step-card">
              <div className="step-num">1</div>
              <div className="step-illu">
                <HedgehogApply />
                <div style={{ position: 'absolute', top: 12, left: 14 }}>
                  <StickerSparkle size={28} tilt={-12} />
                </div>
              </div>
              <h3>Подайте заявку</h3>
              <p>
                Заполните анкету за 3 минуты: название, адрес, тип заведения.
                Мы перезвоним и подтвердим в течение 1–2 рабочих дней.
              </p>
              <div className="step-meta">
                <IconClock size={14} /> 3 минуты на форму
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
              <h3>Соберите сюрприз-бокс</h3>
              <p>
                Положите в коробку то, что осталось на конец дня. Поставьте цену
                втрое ниже обычной — мы покажем бокс клиентам в радиусе 2 км.
              </p>
              <div className="step-meta">
                <IconBox size={14} /> до 5 боксов одновременно
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
              <h3>Получайте выплаты</h3>
              <p>
                Клиент платит в приложении, приходит и забирает по QR-коду.
                Деньги — на счёт каждый понедельник, без бумажек и отчётности.
              </p>
              <div className="step-meta">
                <IconCard size={14} /> выплаты раз в неделю
              </div>
            </div>
          </Reveal>
        </div>
      </div>
    </section>
  )
}

// ─── Benefits ────────────────────────────────────────────────────────────────
function Benefits() {
  return (
    <section className="section benefits">
      <div className="wrap">
        <div className="section-head">
          <SectionMarker num="03" label="Преимущества" />
          <h2 className="h-section">
            Почему партнёры остаются с&nbsp;нами{' '}
            <span style={{ fontFamily: 'var(--f-serif)', fontStyle: 'italic', color: 'var(--c-primary)' }}>
              надолго
            </span>
          </h2>
        </div>

        <div className="benefits-grid">
          <Reveal delay={0}>
            <div className="benefit-card">
              <span className="benefit-tag">Главное</span>
              <div className="benefit-icon"><IconCard size={28} /></div>
              <div className="b-num"><CountUp to={32} suffix="%" /></div>
              <h3>дополнительная выручка с того, что выбрасывалось</h3>
              <p>В среднем партнёр зарабатывает 18&thinsp;000–45&thinsp;000&nbsp;₽/мес на остатках, которые шли в мусор.</p>
            </div>
          </Reveal>

          <Reveal delay={80}>
            <div className="benefit-card">
              <span className="benefit-tag">Маркетинг</span>
              <div className="benefit-icon" style={{ background: 'var(--c-primary-soft)', color: 'var(--c-primary-deep)' }}>
                <IconHandshake size={28} />
              </div>
              <h3>Новые гости, которые приходят за полным меню</h3>
              <p>67% покупателей боксов возвращаются как обычные клиенты — попробовали и понравилось.</p>
            </div>
          </Reveal>

          <Reveal delay={160}>
            <div className="benefit-card">
              <span className="benefit-tag">Льгота</span>
              <div className="benefit-icon" style={{ background: 'var(--c-primary-soft)', color: 'var(--c-primary-deep)' }}>
                <IconShield size={28} />
              </div>
              <h3>10% комиссия первые 3&nbsp;месяца</h3>
              <p>Полная стандартная — 20%. Промо — для всех новых партнёров без условий.</p>
            </div>
          </Reveal>

          <Reveal delay={240}>
            <div className="benefit-card">
              <span className="benefit-tag">Удобно</span>
              <div className="benefit-icon" style={{ background: 'rgba(0,0,0,0.1)', color: 'var(--c-primary-deep)' }}>
                <IconClock size={28} />
              </div>
              <h3>5&nbsp;минут в день на управление</h3>
              <p>Один экран в браузере. Открыли — обновили количество — закрыли.</p>
            </div>
          </Reveal>

          <Reveal delay={320}>
            <div className="benefit-card">
              <span className="benefit-tag">Эко</span>
              <div className="benefit-icon" style={{ background: 'var(--c-primary-soft)', color: 'var(--c-primary-deep)' }}>
                <IconChart size={28} />
              </div>
              <h3>Снижение пищевых отходов на&nbsp;30%</h3>
              <p>Цифра, которую можно показать клиентам, гостям ивента и налоговой.</p>
            </div>
          </Reveal>
        </div>
      </div>
    </section>
  )
}

// ─── Calculator section ───────────────────────────────────────────────────────
function CalcSection() {
  return (
    <section id="calc" className="section calc">
      <div className="wrap">
        <Reveal>
          <LandingCalculator
            SectionMarker={<SectionMarker num="04" label="Калькулятор дохода" />}
          />
        </Reveal>
      </div>
    </section>
  )
}

// ─── Reviews ─────────────────────────────────────────────────────────────────
const REVIEWS = [
  {
    initials: 'АП',
    name: 'Анна Петрова',
    role: 'владелица пекарни «Корица», СПб',
    quote: 'За месяц перестала списывать круассаны. Раньше в 21:00 — две корзины в мусор. Сейчас — две корзины в Бережок и +28 000 ₽ в кассу.',
    metric: '+28k ₽/мес',
    stars: 5,
  },
  {
    initials: 'ИК',
    name: 'Игорь Кравцов',
    role: 'сеть кофеен «Помол», 4 точки',
    quote: 'Думал, будет геморрой с настройкой и обучением баристов. По факту — 10 минут показал девочкам кнопку «Выдать», и всё работает само.',
    metric: '4 точки',
    stars: 5,
  },
  {
    initials: 'МР',
    name: 'Михаил Романов',
    role: 'ресторан «Шумовка», Москва',
    quote: 'Самое неожиданное — половина клиентов Бережка вернулась поужинать обычным меню. Это уже не про остатки, это про привлечение.',
    metric: '67% возврат',
    stars: 5,
  },
]

const PARTNER_NAMES = [
  'Корица', 'Помол', 'Шумовка', 'Хлебный двор', 'Чай & Тосты',
  'Café Saturne', 'Тёплый плед', 'Маркет 24', 'Пироги&Сыр', 'Грильяж',
]

function Reviews() {
  const marqueeNames = [...PARTNER_NAMES, ...PARTNER_NAMES]
  return (
    <section id="reviews" className="section reviews">
      <div className="wrap">
        <div className="section-head">
          <SectionMarker num="05" label="Отзывы партнёров" />
          <h2 className="h-section">
            «Окей, попробуем» —{' '}
            <span style={{ fontFamily: 'var(--f-serif)', fontStyle: 'italic', color: 'var(--c-primary)' }}>
              через&nbsp;неделю
            </span>{' '}
            не отключают
          </h2>
        </div>

        <div className="reviews-grid">
          {REVIEWS.map((r, i) => (
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

      <div className="partner-marquee">
        <div className="partner-marquee-track">
          {marqueeNames.map((n, i) => (
            <span key={i} className="partner-name">★ {n}</span>
          ))}
        </div>
      </div>
    </section>
  )
}

// ─── FAQ ─────────────────────────────────────────────────────────────────────
const FAQS = [
  {
    q: 'А если клиент не пришёл за боксом?',
    a: 'Бокс уже оплачен. Деньги остаются у вас за вычетом комиссии. Еду можете отдать сотрудникам или утилизировать — но вы уже в плюсе.',
  },
  {
    q: 'Сколько боксов в день мне продавать?',
    a: 'Сколько хотите. Лимит — 5 активных боксов одновременно. Большинство пекарен делает 6-10 боксов в день, кофейни — 3-5, рестораны — 2-4.',
  },
  {
    q: 'Что если боюсь за репутацию — клиенты решат, что у нас «несвежее»?',
    a: 'Это не так работает. В боксе — продукция дня, просто не успевшая продаться. По нашей статистике, 67% покупателей боксов возвращаются за обычным меню. Бережок — это не «сток», это окно скидок к закрытию.',
  },
  {
    q: 'Какая комиссия и есть ли скрытые платежи?',
    a: '20% от суммы заказа. Первые 3 месяца — 10% (промо для всех новых партнёров). Для тех, кто делает >100 заказов/мес — 15% (VIP). Других платежей, абонентки или платы за подключение нет.',
  },
  {
    q: 'Когда я получу деньги?',
    a: 'Автоматическая выплата каждый понедельник за предыдущую неделю. На счёт юр.лица или ИП. Реквизиты привязываете один раз в личном кабинете.',
  },
  {
    q: 'Можно интегрировать с моей кассой / iiko / R-Keeper?',
    a: 'Сейчас работаем как отдельная веб-панель — этого хватает большинству. Интеграция с iiko и R-Keeper — в roadmap на 2027 год.',
  },
  {
    q: 'А налоги и чеки?',
    a: 'Мы фискализируем продажу через Yandex Pay/ЮKassa с вашим ОФД. Чек уходит клиенту автоматически. Вам — отчёт за неделю с детализацией.',
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
  const [openIdx, setOpenIdx] = useState(0)
  return (
    <section id="faq" className="section faq">
      <div className="wrap">
        <div className="faq-grid">
          <div className="faq-aside">
            <div style={{ marginBottom: 18 }}>
              <SectionMarker num="06" label="FAQ" />
            </div>
            <h2 className="h-section" style={{ textAlign: 'left' }}>
              Частые вопросы скептика
            </h2>
            <p className="h-sub" style={{ textAlign: 'left' }}>
              Собрали то, что владельцы кафе спрашивают на встречах в первые 10 минут.
              Если что-то не нашли — напишите в форму ниже, ответим лично.
            </p>
            <div style={{ marginTop: 28, display: 'flex', gap: 12 }}>
              <StickerLeaf size={48} tilt={-10} />
              <StickerSparkle size={32} tilt={20} />
            </div>
          </div>
          <div className="faq-list">
            {FAQS.map((f, i) => (
              <FAQItem
                key={i}
                q={f.q}
                a={f.a}
                open={openIdx === i}
                onClick={() => setOpenIdx(openIdx === i ? -1 : i)}
              />
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}

// ─── Final + Form ─────────────────────────────────────────────────────────────
function Final() {
  return (
    <section id="apply" className="final">
      <div className="wrap final-grid">
        <Reveal>
          <div>
            <div style={{ marginBottom: 22 }}>
              <SectionMarker num="07" label="Подключение" dark />
            </div>
            <h2 className="h-section">
              Готовы перестать выбрасывать{' '}
              <span style={{ color: 'var(--c-accent)' }}>деньги в&nbsp;мусор?</span>
            </h2>
            <p className="final-sub" style={{ marginTop: 16 }}>
              Оставьте номер — менеджер свяжется в течение пары часов, поможет
              рассчитать выгоду именно для вашего заведения и подключит за вечер.
            </p>
            <div className="final-checks">
              {[
                'Бесплатное подключение и обучение',
                'Промо-комиссия 10% первые 3 месяца',
                'Без договоров на 100 страниц и обязательств',
                'В любой момент можно отключиться — без штрафов',
              ].map((item) => (
                <div key={item} className="final-check">
                  <span className="final-check-icn"><IconCheck size={14} /></span>
                  <span>{item}</span>
                </div>
              ))}
            </div>
          </div>
        </Reveal>

        <Reveal delay={150}>
          <div className="final-form">
            <h3>Заявка партнёра</h3>
            <LandingFinalForm />
          </div>
        </Reveal>
      </div>
    </section>
  )
}

// ─── Footer ───────────────────────────────────────────────────────────────────
function Footer() {
  return (
    <footer className="footer">
      <div className="footer-brand">
        <HedgehogMark size={28} />
        Бережок
      </div>
      <div>© 2026 Бережок · Меньше отходов — больше добра.</div>
      <div className="footer-links">
        <a href="#">Договор</a>
        <a href="#">Конфиденциальность</a>
        <a href="mailto:support@berezhok.ru">support@berezhok.ru</a>
      </div>
    </footer>
  )
}

// ─── LandingPage ─────────────────────────────────────────────────────────────
export default function LandingPage() {
  return (
    <div className="lp">
      <div className="grain" />
      <Nav />
      <Hero />
      <Steps />
      <Benefits />
      <CalcSection />
      <Reviews />
      <FAQ />
      <Final />
      <Footer />
    </div>
  )
}
