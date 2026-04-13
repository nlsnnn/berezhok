import { useMemo, useState } from "react";
import { observer } from "mobx-react-lite";
import { useNavigate } from "react-router-dom";
import { Building2, Save } from "lucide-react";
import { toast } from "sonner";
import { savePartnerLegalInfo } from "@/api/partner";
import { getErrorMessage } from "@/lib/utils";
import PartnerLayout from "@/components/partner/layout/PartnerLayout";
import Input from "@/components/ui/form/Input";
import Label from "@/components/ui/form/Label";
import Button from "@/components/ui/actions/Button";

const INITIAL_FORM = {
  inn: "",
  ogrn: "",
  kpp: "",
  legal_address: "",
  legal_name: "",
};

function PartnerLegalInfoPageBase() {
  const navigate = useNavigate();
  const [form, setForm] = useState(INITIAL_FORM);
  const [errors, setErrors] = useState({});
  const [submitting, setSubmitting] = useState(false);

  const progress = useMemo(
    () =>
      [
        form.inn,
        form.ogrn,
        form.kpp,
        form.legal_address,
        form.legal_name,
      ].filter(Boolean).length,
    [form.inn, form.ogrn, form.kpp, form.legal_address, form.legal_name],
  );

  const setField = (field) => (e) => {
    const value = e.target.value;
    setForm((prev) => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: "" }));
    }
  };

  const validate = () => {
    const nextErrors = {};

    if (!/^\d{10}$/.test(form.inn)) {
      nextErrors.inn = "ИНН должен содержать 10 цифр";
    }

    if (!/^\d{13}$/.test(form.ogrn)) {
      nextErrors.ogrn = "ОГРН должен содержать 13 цифр";
    }

    if (!/^\d{9}$/.test(form.kpp)) {
      nextErrors.kpp = "КПП должен содержать 9 цифр";
    }

    if (!form.legal_name.trim()) {
      nextErrors.legal_name = "Укажите юридическое наименование";
    }

    if (!form.legal_address.trim()) {
      nextErrors.legal_address = "Укажите юридический адрес";
    }

    return nextErrors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const nextErrors = validate();
    if (Object.keys(nextErrors).length) {
      setErrors(nextErrors);
      return;
    }

    setErrors({});
    setSubmitting(true);

    try {
      await savePartnerLegalInfo({
        inn: form.inn,
        ogrn: form.ogrn,
        kpp: form.kpp,
        legal_address: form.legal_address,
        legal_name: form.legal_name,
      });
      toast.success("Юридические данные сохранены");
      navigate("/partner/dashboard");
    } catch (error) {
      toast.error(getErrorMessage(error));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <PartnerLayout
      title="Юридические данные"
      subtitle="Заполните реквизиты компании, чтобы активировать боксы и начать принимать заказы"
    >
      <div className="max-w-3xl space-y-5">
        <section className="rounded-2xl border border-brand-200 bg-brand-50/60 p-4">
          <div className="flex items-center justify-between gap-3 text-sm">
            <p className="font-medium text-brand-900">
              Прогресс заполнения формы
            </p>
            <p className="text-brand-700">{progress}/5 полей</p>
          </div>
          <div className="mt-2 h-2 overflow-hidden rounded-full bg-brand-100">
            <div
              className="h-full rounded-full bg-brand-500 transition-all"
              style={{ width: `${(progress / 5) * 100}%` }}
            />
          </div>
        </section>

        <form onSubmit={handleSubmit} className="card space-y-5" noValidate>
          <div className="grid sm:grid-cols-2 gap-5">
            <div>
              <Label required>ИНН</Label>
              <Input
                value={form.inn}
                onChange={setField("inn")}
                error={errors.inn}
                inputMode="numeric"
                maxLength={10}
                placeholder="1234567890"
              />
            </div>
            <div>
              <Label required>ОГРН</Label>
              <Input
                value={form.ogrn}
                onChange={setField("ogrn")}
                error={errors.ogrn}
                inputMode="numeric"
                maxLength={13}
                placeholder="1234567890123"
              />
            </div>
          </div>

          <div>
            <Label required>КПП</Label>
            <Input
              value={form.kpp}
              onChange={setField("kpp")}
              error={errors.kpp}
              inputMode="numeric"
              maxLength={9}
              placeholder="123456789"
            />
          </div>

          <div>
            <Label required>Юридическое наименование</Label>
            <Input
              value={form.legal_name}
              onChange={setField("legal_name")}
              error={errors.legal_name}
              placeholder="ООО Вкусная пекарня"
            />
          </div>

          <div>
            <Label required>Юридический адрес</Label>
            <Input
              value={form.legal_address}
              onChange={setField("legal_address")}
              error={errors.legal_address}
              placeholder="г. Москва, ул. Ленина, 1"
            />
          </div>

          <div className="rounded-xl border border-amber-300 bg-amber-50 px-4 py-3 text-sm text-amber-900">
            Пока юридические данные не заполнены, новые боксы можно создавать
            только как неактивные или черновики.
          </div>

          <div className="flex gap-3 pt-1">
            <Button type="submit" className="flex-1" disabled={submitting}>
              {submitting ? (
                "Сохраняем..."
              ) : (
                <>
                  <Save size={16} /> Сохранить реквизиты
                </>
              )}
            </Button>
            <Button
              type="button"
              variant="secondary"
              onClick={() => navigate("/partner/profile")}
            >
              Отмена
            </Button>
          </div>
        </form>

        <div className="rounded-xl border border-cream-200 bg-white px-4 py-3 text-sm text-brand-700 flex items-center gap-2">
          <Building2 size={16} className="text-brand-500" />
          После сохранения данных статус партнера обновится автоматически.
        </div>
      </div>
    </PartnerLayout>
  );
}

export default observer(PartnerLegalInfoPageBase);
