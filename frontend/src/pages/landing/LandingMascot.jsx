function Hedgehog({ pose = 'hi', style = {} }) {
  return (
    <img
      src={`/mascot/${pose}.png`}
      alt=""
      draggable="false"
      style={{
        width: '100%',
        height: '100%',
        objectFit: 'contain',
        userSelect: 'none',
        pointerEvents: 'none',
        display: 'block',
        ...style,
      }}
    />
  )
}

export function HedgehogHero({ pose = 'hi' }) {
  return <Hedgehog pose={pose} />
}

export function HedgehogApply() {
  return <Hedgehog pose="come" />
}

export function HedgehogBox() {
  return <Hedgehog pose="box" />
}

export function HedgehogMoney() {
  return <Hedgehog pose="heart" />
}

export function HedgehogMark({ size = 40 }) {
  return (
    <img
      src="/logo-sm.png"
      alt="Бережок"
      draggable="false"
      style={{
        width: size,
        height: size,
        display: 'block',
        borderRadius: size * 0.22,
        boxShadow: '0 4px 10px -4px rgba(31,60,35,0.35)',
        userSelect: 'none',
      }}
    />
  )
}
