import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import { CheckCircle2, Leaf, ShoppingBag, Handshake, BarChart3, ChevronDown } from 'lucide-react'
import LandingNav from '@/components/LandingNav'
import ApplicationForm from './ApplicationForm'

const HOW_IT_WORKS = [
  {
    icon: Handshake,
    title: 'Подайте заявку',
    desc: 'Заполните короткую форму о вашем заведении. Мы рассмотрим её в течение 1–2 рабочих дней.',
  },
  {
    icon: ShoppingBag,
    title: 'Создайте «коробки»',
    desc: 'Упакуйте остатки дня в сюрприз-боксы по сниженной цене. Максимум 5 активных позиций.',
  },
  {
    icon: BarChart3,
    title: 'Зарабатывайте',
    desc: 'Клиенты находят вас на карте, платят онлайн — вы получаете выплаты каждый понедельник.',
  },
]

const BENEFITS = [
  'Комиссия 10% в первые 3 месяца',
  'Еженедельные автоматические выплаты',
  'Выход на аудиторию студентов и эко-активистов',
  'Снижение пищевых отходов на 30–50%',
  'Личный кабинет и статистика продаж',
]

export default function LandingPage() {
  const [formVisible, setFormVisible] = useState(false)

  return (
    <div className="flex flex-col min-h-screen">
      <LandingNav />

      {/* HERO */}
      <section className="pt-16 bg-gradient-to-br from-cream-50 via-white to-brand-50">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 py-20 sm:py-28 grid lg:grid-cols-2 gap-12 items-center">
          <div>
            <div className="inline-flex items-center gap-2 bg-brand-100 text-brand-700 rounded-full px-4 py-1.5 text-sm font-medium mb-6">
              <Leaf size={14} />
              Фудшеринг-платформа
            </div>
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-brand-900 mb-6 leading-tight">
              Продавайте остатки,
              <br />
              <span className="text-brand-500">спасайте еду</span>
            </h1>
            <p className="text-lg text-brand-600 mb-8 leading-relaxed max-w-lg">
              Бережок помогает кафе, пекарням и ресторанам продавать нераспроданную еду
              в сюрприз-боксах со скидкой 50–70%. Клиенты экономят, вы зарабатываете,
              планета благодарит.
            </p>
            <div className="flex flex-wrap gap-3">
              <a href="#apply" className="btn-primary px-8 py-3 text-base">
                Стать партнёром
              </a>
              <a href="#how-it-works" className="btn-secondary px-8 py-3 text-base">
                Узнать больше
              </a>
            </div>
          </div>

          {/* decorative card */}
          <div className="hidden lg:flex justify-center">
            <div className="relative">
              <div className="w-72 h-72 rounded-3xl bg-brand-500 opacity-10 absolute -top-4 -left-4" />
              <div className="card relative w-72 text-center py-10 shadow-xl">
                <div className="text-6xl mb-4">🥐</div>
                <p className="font-semibold text-brand-800 text-lg">Сюрприз-бокс</p>
                <p className="text-cream-500 text-sm mt-1">Пекарня «Ромашка»</p>
                <div className="mt-4 flex items-center justify-center gap-3">
                  <span className="line-through text-cream-400 text-sm">600 ₽</span>
                  <span className="text-brand-500 font-bold text-2xl">250 ₽</span>
                </div>
                <span className="mt-3 inline-block badge bg-brand-100 text-brand-700">−58%</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* HOW IT WORKS */}
      <section id="how-it-works" className="py-20 bg-white">
        <div className="max-w-6xl mx-auto px-4 sm:px-6">
          <div className="text-center mb-14">
            <h2 className="text-3xl sm:text-4xl font-bold text-brand-900">Как это работает</h2>
            <p className="text-brand-600 mt-3 text-lg">Три шага до первых продаж</p>
          </div>
          <div className="grid md:grid-cols-3 gap-8">
            {HOW_IT_WORKS.map(({ icon: Icon, title, desc }, i) => (
              <div key={i} className="card text-center hover:shadow-md transition-shadow">
                <div className="w-14 h-14 rounded-2xl bg-brand-100 flex items-center justify-center mx-auto mb-5">
                  <Icon size={26} className="text-brand-500" />
                </div>
                <div className="text-xs font-bold text-brand-400 mb-2 uppercase tracking-wider">Шаг {i + 1}</div>
                <h3 className="text-lg font-semibold text-brand-900 mb-3">{title}</h3>
                <p className="text-brand-600 text-sm leading-relaxed">{desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* BENEFITS */}
      <section className="py-20 bg-cream-100">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 grid lg:grid-cols-2 gap-12 items-center">
          <div>
            <h2 className="text-3xl sm:text-4xl font-bold text-brand-900 mb-6">
              Почему партнёры выбирают Бережок
            </h2>
            <ul className="space-y-4">
              {BENEFITS.map((b, i) => (
                <li key={i} className="flex items-start gap-3">
                  <CheckCircle2 size={20} className="text-brand-500 mt-0.5 shrink-0" />
                  <span className="text-brand-700">{b}</span>
                </li>
              ))}
            </ul>
          </div>
          <div className="card bg-white shadow-xl">
            <div className="text-center mb-6">
              <div className="text-4xl font-bold text-brand-500">30%</div>
              <div className="text-brand-600 text-sm mt-1">среднее снижение пищевых отходов</div>
            </div>
            <div className="grid grid-cols-2 gap-4 text-center">
              {[
                { value: '500+', label: 'партнёров в России' },
                { value: '50%', label: 'средняя скидка для клиентов' },
                { value: '10%', label: 'комиссия первые 3 месяца' },
                { value: '1–2', label: 'дня рассмотрения заявки' },
              ].map(({ value, label }, i) => (
                <div key={i} className="bg-cream-50 rounded-xl p-4">
                  <div className="text-2xl font-bold text-brand-700">{value}</div>
                  <div className="text-xs text-brand-500 mt-1 leading-snug">{label}</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* APPLICATION FORM */}
      <section id="apply" className="py-20 bg-white">
        <div className="max-w-2xl mx-auto px-4 sm:px-6">
          <div className="text-center mb-10">
            <h2 className="text-3xl sm:text-4xl font-bold text-brand-900">Подать заявку</h2>
            <p className="text-brand-600 mt-3">Заполните форму, и мы свяжемся с вами в течение 1–2 дней</p>
          </div>
          <ApplicationForm />
        </div>
      </section>

      {/* FOOTER */}
      <footer className="bg-brand-900 text-brand-300 py-10 mt-auto">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2 text-white font-semibold">
            <Leaf size={18} className="text-brand-400" />
            Бережок
          </div>
          <p className="text-sm text-center">© 2026 Бережок. Меньше отходов — больше добра.</p>
        </div>
      </footer>
    </div>
  )
}
