# sport-hub-profile
user profile

โครงที่แนะนำ
1. plans

เก็บ master ของแพ็กเกจ
CREATE TABLE public.plans (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    code varchar(50) NOT NULL,
    name varchar(100) NOT NULL,
    description text NULL,
    price numeric(10,2) NOT NULL DEFAULT 0,
    billing_cycle varchar(20) NOT NULL DEFAULT 'monthly',
    trial_days int NOT NULL DEFAULT 0,

    max_fields int NULL,
    max_courts int NULL,
    has_reports boolean NOT NULL DEFAULT false,
    has_promotion boolean NOT NULL DEFAULT false,
    has_homepage_feature boolean NOT NULL DEFAULT false,
    has_priority_support boolean NOT NULL DEFAULT false,

    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,

    CONSTRAINT plans_pkey PRIMARY KEY (id),
    CONSTRAINT plans_code_key UNIQUE (code)
);

INSERT INTO public.plans
(code, name, description, price, billing_cycle, trial_days, max_fields, max_courts, has_reports, has_promotion, has_homepage_feature, has_priority_support)
VALUES
('free', 'Free', 'แพ็กเกจเริ่มต้น', 0, 'monthly', 0, 1, 2, false, false, false, false),
('standard', 'Standard', 'แพ็กเกจสำหรับสนามที่ต้องการเติบโต', 299, 'monthly', 0, 3, NULL, true, true, false, false),
('pro', 'Pro', 'แพ็กเกจเต็มรูปแบบ', 999, 'monthly', 7, NULL, NULL, true, true, true, true);

2. subscriptions
เก็บว่าผู้ใช้คนนี้กำลังใช้ plan อะไรอยู่ และมีประวัติย้อนหลัง

CREATE TABLE public.subscriptions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    plan_id uuid NOT NULL,

    status varchar(30) NOT NULL,
    -- trial / active / expired / cancelled

    start_at timestamp NOT NULL,
    end_at timestamp NOT NULL,

    trial_start_at timestamp NULL,
    trial_end_at timestamp NULL,

    activated_at timestamp NULL,
    expired_at timestamp NULL,
    cancelled_at timestamp NULL,
    cancel_reason text NULL,

    auto_renew boolean NOT NULL DEFAULT false,

    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,

    CONSTRAINT subscriptions_pkey PRIMARY KEY (id),
    CONSTRAINT subscriptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id),
    CONSTRAINT subscriptions_plan_id_fkey FOREIGN KEY (plan_id) REFERENCES public.plans(id)
);

subscription status

ใช้แค่ 4 ค่า:
- trial
- active
- expired
- cancelled