import { useMemo, useState, useEffect } from "react";
import { observer } from "mobx-react-lite";
import { useNavigate } from "react-router-dom";
import { Building2, Save, CheckCircle2, Clock, XCircle, AlertTriangle } from "lucide-react";
import { toast } from "sonner";
import { savePartnerLegalInfo } from "@/api/partner";
import { getErrorMessage, getValidationDetails } from "@/lib/utils";
import { useStores } from "@/context/StoresContext";
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

const SERVER_FIELD_MESSAGES = {
  inn: {
    fallback: "Проверьте ИНН",
    required: "Укажите ИНН",
    numeric: "ИНН должен содержать только цифры",
    min: "ИНН должен содержать 10 или 12 цифр",
    max: "ИНН должен содержать 10 или 12 цифр",
  },
  ogrn: {
    fallback: "Проверьте ОГРН",
    numeric: "ОГРН должен содержать только цифры",
    min: "ОГРН должен содержать 13 или 15 цифр",
    max: "ОГРН должен содержать 13 или 15 цифр",
  },
  kpp: {
    fallback: "Проверьте КПП",
    numeric: "КПП должен содержать только цифры",
    len: "КПП должен содержать 9 цифр",
  },
  legal_address: {
    fallback: "Проверьте юридический адрес",
    required: "Укажите юридический адрес",
    min: "Юридический адрес должен содержать минимум 5 символов",
    max: "Юридический адрес должен содержать не больше 500 символов",
  },
  legal_name: {
    fallback: "Проверьте юридическое наименование",
    required: "Укажите юридическое наименование",
    min: "Юридическое наименование должно содержать минимум 2 символа",
    max: "Юридическое наименование должно содержать не больше 200 символов",
  },
};

function translateServerFieldError(field, message) {
  const messages = SERVER_FIELD_MESSAGES[field];
  if (!messages) return String(message || "Некорректное значение");

  const normalized = String(message || "").toLowerCase();
  if (normalized.includes("required")) return messages.required || messages.fallback;
  if (normalized.includes("numeric")) return messages.numeric || messages.fallback;
  if (normalized.includes("at least") || normalized.includes("min")) return messages.min || messages.fallback;
  if (normalized.includes("at most") || normalized.includes("max")) return messages.max || messages.fallback;
  if (normalized.includes("length") || normalized.includes("len")) return messages.len || messages.fallback;

  return messages.fallback;
}

function normalizeServerValidationDetails(details) {
  return Object.fromEntries(
    Object.entries(details).map(([field, message]) => [
      field,
      translateServerFieldError(field, message),
    ]),
  );
}

function fieldErrorFromServerMessage(message) {
  const normalized = String(message || "").toLowerCase();
  if (normalized.includes("invalid inn")) {
    return { inn: "Проверьте ИНН: контрольная сумма не совпадает" };
  }
  if (normalized.includes("invalid ogrn")) {
    return { ogrn: "Проверьте формат ОГРН" };
  }
  if (normalized.includes("invalid kpp")) {
    return { kpp: "КПП должен содержать 9 цифр" };
  }
  return null;
}

function PartnerLegalInfoPageBase() {
  const navigate = useNavigate();
  const { dashboardStore } = useStores();
  const [form, setForm] = useState(INITIAL_FORM);
  const [errors, setErrors] = useState({});
  const [submitting, setSubmitting] = useState(false);
  const [formUnlocked, setFormUnlocked] = useState(false);

  useEffect(() => {
    if (!dashboardStore.data) dashboardStore.load();
  }, [dashboardStore]);

  const hasLegalInfo = dashboardStore.data?.partner?.has_legal_info;
  const legalInfoStatus = dashboardStore.data?.partner?.legal_info_status;

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

    if (!/^(\d{10}|\d{12})$/.test(form.inn)) {
      nextErrors.inn = "ИНН должен содержать 10 или 12 цифр";
    }

    if (!/^(\d{13}|\d{15})$/.test(form.ogrn)) {
      nextErrors.ogrn = "ОГРН должен содержать 13 или 15 цифр";
    }

    if (!/^\d{9}$/.test(form.kpp)) {
      nextErrors.kpp = "КПП должен содержать 9 цифр";
    }

    if (form.legal_name.trim().length < 2) {
      nextErrors.legal_name = "Юридическое наименование должно содержать минимум 2 символа";
    }

    if (form.legal_address.trim().length < 5) {
      nextErrors.legal_address = "Юридический адрес должен содержать минимум 5 символов";
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
      const validationDetails = getValidationDetails(error);
      if (validationDetails) {
        setErrors(normalizeServerValidationDetails(validationDetails));
        toast.error("Проверьте выделенные поля");
        return;
      }

      const message = getErrorMessage(error);
      const fieldErrors = fieldErrorFromServerMessage(message);
      if (fieldErrors) {
        setErrors(fieldErrors);
        toast.error("Проверьте выделенные поля");
        return;
      }

      toast.error(message);
    } finally {
      setSubmitting(false);
    }
  };

  const layoutProps = {
    title: "Юридические данные",
    subtitle: "Заполните реквизиты компании, чтобы активировать боксы и начать принимать заказы",
  };

  if (dashboardStore.loading && !dashboardStore.data) {
    return (
      <PartnerLayout {...layoutProps}>
        <div className="max-w-3xl">
          <div className="h-32 flex items-center justify-center text-sm text-brand-400">Загрузка...</div>
        </div>
      </PartnerLayout>
    );
  }

  if (legalInfoStatus === "verified") {
    return (
      <PartnerLayout {...layoutProps}>
        <div className="max-w-3xl">
          <div className="rounded-2xl border border-green-200 bg-gradient-to-r from-green-50 to-emerald-50 p-5 md:p-6">
            <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
              <div className="space-y-2">
                <div className="inline-flex items-center gap-2 rounded-full border border-green-200 bg-white/70 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-green-700">
                  <CheckCircle2 size={14} />
                  Данные подтверждены
                </div>
                <h2 className="text-lg font-bold text-green-950">Юридические данные приняты</h2>
                <p className="text-sm text-green-900/90">
                  Ваши данные проверены администратором — аккаунт полностью активен и готов к работе.
                </p>
              </div>
            </div>
            <div className="mt-5 pt-4 border-t border-green-200/60">
              <Button variant="secondary" onClick={() => navigate("/partner/profile")}>
                Вернуться в профиль
              </Button>
            </div>
          </div>
        </div>
      </PartnerLayout>
    );
  }

  if (hasLegalInfo && legalInfoStatus === "pending" && !formUnlocked) {
    return (
      <PartnerLayout {...layoutProps}>
        <div className="max-w-3xl">
          <div className="rounded-2xl border border-blue-200 bg-gradient-to-r from-blue-50 to-indigo-50 p-5 md:p-6">
            <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
              <div className="space-y-2">
                <div className="inline-flex items-center gap-2 rounded-full border border-blue-200 bg-white/70 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-blue-700">
                  <Clock size={14} />
                  На проверке
                </div>
                <h2 className="text-lg font-bold text-blue-950">Данные отправлены — ожидайте проверки</h2>
                <p className="text-sm text-blue-900/90">
                  Администратор проверит ваши юридические данные и активирует аккаунт.
                  Это обычно занимает до 1 рабочего дня.
                </p>
              </div>
              <div className="shrink-0 rounded-xl border border-blue-200 bg-white/70 px-4 py-3 text-center">
                <p className="text-xs text-blue-600 font-medium">Статус</p>
                <p className="text-sm font-bold text-blue-800 mt-0.5">На рассмотрении</p>
              </div>
            </div>
            <div className="mt-5 pt-4 border-t border-blue-200/60 flex flex-wrap gap-3">
              <Button variant="secondary" onClick={() => navigate("/partner/profile")}>
                Вернуться в профиль
              </Button>
              <button
                type="button"
                onClick={() => setFormUnlocked(true)}
                className="text-sm text-blue-600 hover:text-blue-800 underline underline-offset-2 transition-colors self-center"
              >
                Изменить данные
              </button>
            </div>
          </div>
        </div>
      </PartnerLayout>
    );
  }

  return (
    <PartnerLayout {...layoutProps}>
      <div className="max-w-3xl space-y-5">
        {legalInfoStatus === "failed" && (
          <div className="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 flex items-start gap-3 text-sm text-red-900">
            <XCircle size={16} className="text-red-500 shrink-0 mt-0.5" />
            <p>
              <span className="font-semibold">Данные отклонены.</span>{" "}
              Администратор отклонил ваши данные. Проверьте корректность ИНН, ОГРН и юридического адреса.
            </p>
          </div>
        )}

        {formUnlocked && legalInfoStatus === "pending" && (
          <div className="rounded-2xl border border-amber-300 bg-amber-50 px-4 py-3 flex items-start gap-3 text-sm text-amber-900">
            <AlertTriangle size={16} className="text-amber-500 shrink-0 mt-0.5" />
            <p>
              <span className="font-semibold">Данные уже на проверке.</span>{" "}
              Повторная отправка заменит текущую заявку и сбросит статус проверки.
            </p>
          </div>
        )}

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
                maxLength={12}
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
                maxLength={15}
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
          После сохранения данные уйдут администратору на проверку.
        </div>
      </div>
    </PartnerLayout>
  );
}

export default observer(PartnerLegalInfoPageBase);
