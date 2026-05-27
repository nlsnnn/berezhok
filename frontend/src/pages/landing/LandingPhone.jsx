// LandingPhone.jsx — iPhone-style mockup components for the consumer home page.
// Uses CSS classes from landing.css (.lp .phone-*)

function PhoneStatusBar({ time = '9:41' }) {
  return (
    <div className="phone-statusbar">
      <span>{time}</span>
      <div className="phone-statusbar__icons">
        {/* signal */}
        <svg width="16" height="10" viewBox="0 0 16 10" fill="currentColor">
          <rect x="0"  y="6" width="2.5" height="4" rx="0.6" />
          <rect x="4"  y="4" width="2.5" height="6" rx="0.6" />
          <rect x="8"  y="2" width="2.5" height="8" rx="0.6" />
          <rect x="12" y="0" width="2.5" height="10" rx="0.6" />
        </svg>
        {/* wifi */}
        <svg width="14" height="10" viewBox="0 0 14 10" fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round">
          <path d="M 1 4 a 8 8 0 0 1 12 0" />
          <path d="M 3 6 a 5 5 0 0 1 8 0" />
          <circle cx="7" cy="9" r="0.9" fill="currentColor" stroke="none" />
        </svg>
        {/* battery */}
        <svg width="22" height="10" viewBox="0 0 22 10" fill="none" stroke="currentColor" strokeWidth="1">
          <rect x="0.5" y="0.5" width="18" height="9" rx="2.5" />
          <rect x="20" y="3" width="1.4" height="4" rx="0.4" fill="currentColor" stroke="none" />
          <rect x="2" y="2" width="14" height="6" rx="1.4" fill="currentColor" stroke="none" />
        </svg>
      </div>
    </div>
  )
}

function PhoneTabs({ active = 'map' }) {
  const tabs = [
    {
      id: 'map', label: 'Карта',
      icn: (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
          <path d="M12 21s-7-7.5-7-12a7 7 0 1 1 14 0c0 4.5-7 12-7 12z" />
          <circle cx="12" cy="9" r="2.5" />
        </svg>
      ),
    },
    {
      id: 'list', label: 'Список',
      icn: (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round">
          <line x1="4" y1="6" x2="20" y2="6" />
          <line x1="4" y1="12" x2="20" y2="12" />
          <line x1="4" y1="18" x2="20" y2="18" />
        </svg>
      ),
    },
    {
      id: 'orders', label: 'Заказы',
      icn: (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
          <path d="M4 7l8-3 8 3v10l-8 3-8-3z" />
          <path d="M4 7l8 3 8-3" />
          <path d="M12 10v10" />
        </svg>
      ),
    },
    {
      id: 'me', label: 'Я',
      icn: (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round">
          <circle cx="12" cy="8" r="4" />
          <path d="M4 21a8 8 0 0 1 16 0" />
        </svg>
      ),
    },
  ]
  return (
    <div className="phone-tabs">
      {tabs.map((t) => (
        <div key={t.id} className={`phone-tabs__item${active === t.id ? ' is-active' : ''}`}>
          <div className="phone-tabs__icn">{t.icn}</div>
          {t.label}
        </div>
      ))}
    </div>
  )
}

function MapBackground() {
  return (
    <>
      <svg className="phone-map__blocks" viewBox="0 0 320 600" preserveAspectRatio="none">
        <rect x="20"  y="220" width="80"  height="50" rx="6" />
        <rect x="130" y="180" width="60"  height="40" rx="6" />
        <rect x="210" y="230" width="80"  height="60" rx="6" />
        <rect x="30"  y="320" width="70"  height="80" rx="6" />
        <rect x="140" y="380" width="100" height="60" rx="6" />
        <rect x="260" y="370" width="50"  height="40" rx="6" />
        <rect x="20"  y="450" width="60"  height="50" rx="6" />
        <rect x="110" y="490" width="80"  height="40" rx="6" />
        <rect x="220" y="470" width="80"  height="60" rx="6" />
      </svg>
      <svg className="phone-map__streets" viewBox="0 0 320 600" preserveAspectRatio="none">
        <line x1="0"   y1="180" x2="320" y2="195" strokeWidth="14" />
        <line x1="0"   y1="300" x2="320" y2="318" strokeWidth="10" />
        <line x1="0"   y1="445" x2="320" y2="455" strokeWidth="12" />
        <line x1="110" y1="0"   x2="120" y2="600" strokeWidth="14" />
        <line x1="200" y1="0"   x2="215" y2="600" strokeWidth="10" />
      </svg>
    </>
  )
}

export function PhoneMap() {
  return (
    <div className="phone">
      <div className="phone-screen">
        <PhoneStatusBar />
        <div className="phone-map">
          <MapBackground />

          <div className="phone-search">
            <span className="phone-search__icn">
              <svg width="10" height="10" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="2.4" strokeLinecap="round">
                <circle cx="7" cy="7" r="4.5" />
                <line x1="10.5" y1="10.5" x2="14" y2="14" />
              </svg>
            </span>
            <span className="phone-search__txt">Поиск рядом • Москва</span>
          </div>

          <div className="phone-chips">
            <span className="phone-chip is-active">Все</span>
            <span className="phone-chip">Пекарни</span>
            <span className="phone-chip">Кофе</span>
            <span className="phone-chip">До 300₽</span>
          </div>

          <div className="phone-here" style={{ left: '44%', top: '52%' }} />

          <div className="phone-pin phone-pin--accent" style={{ left: '16%', top: '30%', animationDelay: '0s' }}>
            <span className="phone-pin__price">190₽</span>
          </div>
          <div className="phone-pin" style={{ left: '62%', top: '27%', animationDelay: '0.18s' }}>
            <span className="phone-pin__price">240₽</span>
          </div>
          <div className="phone-pin" style={{ left: '30%', top: '60%', animationDelay: '0.36s' }}>
            <span className="phone-pin__price">120₽</span>
          </div>
          <div className="phone-pin phone-pin--accent" style={{ left: '70%', top: '62%', animationDelay: '0.54s' }}>
            <span className="phone-pin__price">350₽</span>
          </div>
          <div className="phone-pin" style={{ left: '20%', top: '72%', animationDelay: '0.72s' }}>
            <span className="phone-pin__price">180₽</span>
          </div>

          <div className="phone-sheet">
            <div className="phone-sheet__pic">🥐</div>
            <div>
              <h4 className="phone-sheet__title">Утренний бокс</h4>
              <div className="phone-sheet__meta">
                <span>пекарня «Корица»</span>
                <span>•</span>
                <span>180 м</span>
              </div>
            </div>
            <div className="phone-sheet__price">
              190₽
              <span className="phone-sheet__strike">590₽</span>
            </div>
          </div>
        </div>
        <PhoneTabs active="map" />
      </div>
    </div>
  )
}

export function PhoneBoxDetail() {
  return (
    <div className="phone">
      <div className="phone-screen">
        <PhoneStatusBar />
        <div className="phone-box">
          <div className="phone-box__hero">
            <span className="phone-box__back">‹</span>
            <span className="phone-box__badge">−68%</span>
            <div style={{ fontSize: 72, position: 'relative', zIndex: 1 }}>🥐</div>
          </div>
          <div className="phone-box__body">
            <h3 className="phone-box__title">Утренний бокс «Корица»</h3>
            <div className="phone-box__sub">
              <b>★ 4.9</b> · 312 отзывов · 180 м от вас
            </div>
            <div className="phone-box__row">
              <span>Забрать</span>
              <b>сегодня, 18:00–19:00</b>
            </div>
            <div className="phone-box__row">
              <span>Внутри</span>
              <b>2–3 круассана, булочки</b>
            </div>
            <div className="phone-box__row">
              <span>Осталось</span>
              <b style={{ color: 'var(--c-warn)' }}>2 бокса</b>
            </div>
            <div className="phone-box__price">
              <b>190₽</b>
              <s>590₽</s>
            </div>
            <div className="phone-box__btn">Забронировать →</div>
          </div>
        </div>
        <PhoneTabs active="list" />
      </div>
    </div>
  )
}

export function PhoneOrdersScreen() {
  return (
    <div className="phone">
      <div className="phone-screen">
        <PhoneStatusBar />
        <div className="phone-orders">
          <h3 className="phone-orders__title">Мои заказы</h3>

          <div className="phone-order phone-order--active">
            <div className="phone-order__icn">🥐</div>
            <div>
              <h4 className="phone-order__t">Утренний бокс «Корица»</h4>
              <div className="phone-order__m">Сегодня · 18:00–19:00</div>
            </div>
            <span className="phone-order__status">Готов</span>
          </div>

          <div className="phone-qr">
            <div className="phone-qr__code">
              <svg viewBox="0 0 50 50">
                <rect width="50" height="50" fill="#fff" />
                {[
                  [2,2,10,10],[38,2,10,10],[2,38,10,10],
                  [4,4,6,6],[40,4,6,6],[4,40,6,6],
                  [14,2,2,2],[18,2,4,2],[24,2,2,4],[28,2,2,2],[32,2,4,2],
                  [2,14,2,2],[2,18,2,4],[2,24,4,2],[2,30,2,2],[2,34,2,2],
                  [14,14,4,4],[20,14,2,2],[24,14,2,4],[28,14,4,2],[34,14,2,4],
                  [14,20,2,4],[18,22,4,2],[24,20,4,4],[30,20,2,2],[34,22,4,2],
                  [14,28,4,2],[20,28,2,4],[26,28,4,2],[32,28,2,4],[38,28,2,2],
                  [14,34,2,4],[18,36,4,2],[24,34,2,2],[28,34,4,4],[34,34,2,4],[38,38,4,4],
                  [22,40,2,2],[26,40,2,4],[32,42,2,2],
                ].map(([x, y, w, h], i) => (
                  <rect key={i} x={x} y={y} width={w} height={h} fill="#1f3c23" />
                ))}
              </svg>
            </div>
            <div>
              <div className="phone-qr__t">Покажите на выдаче</div>
              <div className="phone-qr__code-num">AB12CD34</div>
            </div>
          </div>

          <div className="phone-order">
            <div className="phone-order__icn" style={{ background: 'var(--c-primary-soft)' }}>☕</div>
            <div>
              <h4 className="phone-order__t">Вечерний кофе-бокс «Помол»</h4>
              <div className="phone-order__m">Вчера · получен</div>
            </div>
            <span className="phone-order__status" style={{ background: 'var(--sand-deep)', color: 'var(--c-ink)' }}>Завершён</span>
          </div>

          <div className="phone-order">
            <div className="phone-order__icn" style={{ background: 'var(--c-primary-soft)' }}>🍱</div>
            <div>
              <h4 className="phone-order__t">Обеденный бокс «Шумовка»</h4>
              <div className="phone-order__m">3 дня назад · получен</div>
            </div>
            <span className="phone-order__status" style={{ background: 'var(--sand-deep)', color: 'var(--c-ink)' }}>Завершён</span>
          </div>
        </div>
        <PhoneTabs active="orders" />
      </div>
    </div>
  )
}
