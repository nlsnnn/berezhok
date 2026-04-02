import { CheckCircle2, Handshake, Leaf, ShoppingBag, BarChart3 } from 'lucide-react'
import LandingNav from '@/components/LandingNav'
import ApplicationForm from './ApplicationForm'

const HOW_IT_WORKS = [
  { icon: Handshake, title: 'Подайте заявку', desc: 'Заполните форму о заведении. Ответ обычно в течение 1-2 рабочих дней.' },
  { icon: ShoppingBag, title: 'Создайте боксы', desc: 'Соберите сюрприз-боксы из остатков дня и опубликуйте в приложении.' },
  { icon: BarChart3, title: 'Зарабатывайте', desc: 'Принимайте заказы, подтверждайте выдачу и получайте выплаты.' },
]

const BENEFITS = [
  'Комиссия 10% в первые 3 месяца',
  'Еженедельные автоматические выплаты',
  'Рост повторных заказов и узнаваемости',
  'Снижение пищевых отходов',
  'Личный кабинет с метриками',
]

export default function LandingPage() {
  return (
    <div className="flex flex-col min-h-screen">
      <LandingNav />

      <section className="pt-16 bg-gradient-to-br from-cream-50 via-white to-brand-50">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 py-20 sm:py-28 grid lg:grid-cols-2 gap-12 items-center">
          <div>
            <div className="inline-flex items-center gap-2 bg-brand-100 text-brand-700 rounded-full px-4 py-1.5 text-sm font-medium mb-6">
              <Leaf size={14} /> Фудшеринг-платформа
            </div>
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-brand-900 mb-6 leading-tight">
              Продавайте остатки, <span className="text-brand-500">спасайте еду</span>
            </h1>
            <p className="text-lg text-brand-600 mb-8 leading-relaxed max-w-lg">
              Бережок помогает кафе, пекарням и ресторанам продавать еду в сюрприз-боксах со скидкой.
            </p>
            <div className="flex flex-wrap gap-3">
              <a href="#apply" className="btn-primary px-8 py-3 text-base">Стать партнером</a>
              <a href="#how-it-works" className="btn-secondary px-8 py-3 text-base">Как это работает</a>
            </div>
          </div>

          <div className="hidden lg:flex justify-center">
            <div className="relative">
              <div className="w-72 h-72 rounded-3xl bg-brand-500 opacity-10 absolute -top-4 -left-4" />
              <div className="card relative w-72 text-center py-10 shadow-xl">
                <img src="/logo.png" alt="logo" className="w-20 h-20 rounded-2xl object-cover mx-auto mb-4" />
                <p className="font-semibold text-brand-800 text-lg">Сюрприз-бокс</p>
                <p className="text-cream-500 text-sm mt-1">Экономия и меньше отходов</p>
                <div className="mt-4 flex items-center justify-center gap-3">
                  <span className="line-through text-cream-400 text-sm">600 ₽</span>
                  <span className="text-brand-500 font-bold text-2xl">250 ₽</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="how-it-works" className="py-20 bg-white">
        <div className="max-w-6xl mx-auto px-4 sm:px-6">
          <div className="text-center mb-14">
            <h2 className="text-3xl sm:text-4xl font-bold text-brand-900">Как это работает</h2>
            <p className="text-brand-600 mt-3 text-lg">Три шага до первых заказов</p>
          </div>
          <div className="grid md:grid-cols-3 gap-8">
            {HOW_IT_WORKS.map(({ icon: Icon, title, desc }, i) => (
              <div key={i} className="card text-center hover:shadow-md transition-shadow">
                <div className="w-14 h-14 rounded-2xl bg-brand-100 flex items-center justify-center mx-auto mb-5">
                  <Icon size={26} className="text-brand-500" />
                </div>
                <h3 className="text-lg font-semibold text-brand-900 mb-3">{title}</h3>
                <p className="text-brand-600 text-sm leading-relaxed">{desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="py-20 bg-cream-100">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 grid lg:grid-cols-2 gap-12 items-center">
          <div>
            <h2 className="text-3xl sm:text-4xl font-bold text-brand-900 mb-6">Почему партнеры выбирают Бережок</h2>
            <ul className="space-y-4">
              {BENEFITS.map((item) => (
                <li key={item} className="flex items-start gap-3">
                  <CheckCircle2 size={20} className="text-brand-500 mt-0.5 shrink-0" />
                  <span className="text-brand-700">{item}</span>
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
                { value: '500+', label: 'партнеров' },
                { value: '50%', label: 'средняя скидка' },
                { value: '10%', label: 'комиссия в промо' },
                { value: '1-2', label: 'дня на проверку' },
              ].map(({ value, label }) => (
                <div key={label} className="bg-cream-50 rounded-xl p-4">
                  <div className="text-2xl font-bold text-brand-700">{value}</div>
                  <div className="text-xs text-brand-500 mt-1">{label}</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      <section id="apply" className="py-20 bg-white">
        <div className="max-w-2xl mx-auto px-4 sm:px-6">
          <div className="text-center mb-10">
            <h2 className="text-3xl sm:text-4xl font-bold text-brand-900">Подать заявку</h2>
            <p className="text-brand-600 mt-3">Заполните форму, и мы свяжемся с вами в ближайшее время</p>
          </div>
          <ApplicationForm />
        </div>
      </section>

      <footer className="bg-brand-900 text-brand-300 py-10 mt-auto">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2 text-white font-semibold">
            <img src="/logo.png" alt="Бережок" className="w-6 h-6 rounded-md object-cover" />
            Бережок
          </div>
          <p className="text-sm text-center">© 2026 Бережок. Меньше отходов - больше добра.</p>
        </div>
      </footer>
    </div>
  )
}
