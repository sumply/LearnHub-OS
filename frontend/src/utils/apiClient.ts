// API клиент для работы с бэкендом
// Базовый URL API
// В режиме разработки используем прокси Vite для обхода CORS
const API_BASE_URL =
  import.meta.env.VITE_API_URL || (import.meta.env.DEV ? '/api' : 'http://185.152.92.245:8000');

// Флаг для работы без бэкенда (локальная разработка) - отключен
const LOCAL_DEV_MODE = false;

// Типы согласно документации API
// Роли теперь строковые: "student", "teacher", "admin"
export type UserRole = 'student' | 'teacher' | 'admin' | 'root';

export type ProgressStatus = 0 | 1 | 2 | 3;
export const ProgressStatus = {
  NOT_STARTED: 0,
  IN_PROGRESS: 1,
  AWAITING_REVIEW: 2,
  COMPLETED: 3,
} as const;

export type AnswerStatus = 0 | 1 | 2 | 3;
export const AnswerStatus = {
  NOT_ANSWERED: 0,
  AWAITING_REVIEW: 1,
  INCORRECT: 2,
  CORRECT: 3,
} as const;

// Типы запросов и ответов
export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  refresh_token: string;
  access_token: string;
}

/** Фактические поля ответа POST /auth/login (разные версии бэкенда) */
export interface LoginApiResponse {
  id?: string;
  user_id?: string;
  role?: string;
  access_token?: string;
  /** Частый алиас access_token в ответах API */
  token?: string;
  accessToken?: string;
  refresh_token?: string;
  /** Наш бэкенд: { jwt: { access_token, refresh_token } | string, user: {...} } */
  jwt?: string | Record<string, unknown>;
  user?: { id?: string | number; role?: string };
}

function pickString(...vals: Array<unknown>): string | null {
  for (const v of vals) {
    if (typeof v === 'string' && v.trim() !== '') return v.trim();
  }
  return null;
}

/** Достаёт access token из тела ответа логина / refresh (вложенные форматы учитываются). */
export function extractAccessTokenFromPayload(data: unknown): string | null {
  if (typeof data === 'string' && data.trim() !== '') return data.trim();
  if (!data || typeof data !== 'object') return null;
  const d = data as Record<string, unknown>;

  const jwtVal = d.jwt;
  if (typeof jwtVal === 'string' && jwtVal.trim() !== '') return jwtVal.trim();
  if (jwtVal && typeof jwtVal === 'object') {
    const inner = extractAccessTokenFromPayload(jwtVal);
    if (inner) return inner;
  }

  const nested = d.data;
  if (nested && typeof nested === 'object') {
    const fromData = extractAccessTokenFromPayload(nested);
    if (fromData) return fromData;
  }

  return pickString(d.access_token, d.token, d.accessToken, d.access);
}

function extractRefreshTokenFromPayload(data: unknown): string | null {
  if (!data || typeof data !== 'object') return null;
  const d = data as Record<string, unknown>;

  const jwtVal = d.jwt;
  if (jwtVal && typeof jwtVal === 'object') {
    const inner = extractRefreshTokenFromPayload(jwtVal);
    if (inner) return inner;
  }

  const nested = d.data;
  if (nested && typeof nested === 'object') {
    const fromData = extractRefreshTokenFromPayload(nested);
    if (fromData) return fromData;
  }

  return pickString(d.refresh_token, d.refreshToken, d.refresh);
}

function isPublicAuthPath(endpoint: string): boolean {
  return endpoint === '/auth' || endpoint.startsWith('/auth/');
}

function readJwtPayload(accessToken: string): { sub?: string; user_id?: string; role?: string } | null {
  try {
    const part = accessToken.split('.')[1];
    if (!part) return null;
    const base64 = part.replace(/-/g, '+').replace(/_/g, '/');
    const json = atob(base64);
    return JSON.parse(json) as { sub?: string; user_id?: string; role?: string };
  } catch {
    return null;
  }
}

export interface UserShort {
  id: number;
  short_name: string;
}

export interface UserFull {
  id: string | number;
  first_name: string;
  last_name: string;
  middle_name?: string;
  role: UserRole | string;
}

export interface UserCreateRequest {
  first_name: string;
  last_name: string;
  middle_name?: string;
  email: string;
  role: UserRole | string;
  subject_ids?: number[]; // Для учителей - предметы
  group_ids?: number[]; // Для учителей - группы
}

export interface GroupResponse {
  id: string; // UUID
  name: string;
  curator: UserShort | null; // Может быть null если куратор не назначен
  students: UserShort[]; // Массив студентов группы
}

export interface GroupCreateRequest {
  name: string;
  curator_id?: string; // UUID, необязательно
  student_ids?: string[]; // UUID[], необязательно
}

export interface GroupAddStudentsRequest {
  student_ids: number[];
}

export interface SubjectResponse {
  id: number;
  name: string;
}

export interface SubjectCreateRequest {
  name: string;
}

// Новый формат вопроса для бэкенда
export interface QuizQuestionDetailsSingle {
  options: string[];
  correct: string; // Одна строка для single choice
}

export interface QuizQuestionDetailsMultiple {
  options: string[];
  correct: string[]; // Массив строк для multiple choice
}

export interface QuizQuestionDetailsNumeric {
  correct: number; // Число для numeric
}

export type QuizQuestionDetails = QuizQuestionDetailsSingle | QuizQuestionDetailsMultiple | QuizQuestionDetailsNumeric;

export interface QuizQuestion {
  text: string;
  score: number; // Баллы за вопрос
  type: 'single' | 'multiple' | 'numeric';
  details: QuizQuestionDetails;
}

export interface QuizCreateRequest {
  title: string;
  owner_id: string; // UUID владельца
  summary: string;
  subject_id: string; // UUID предмета или "все темы"
  deadline?: string; // ISO дата/время (необязательно)
  group_ids?: string[]; // UUID[] групп или ["общий"]
  max_attempts?: number; // Максимальное количество попыток (необязательно)
  questions: QuizQuestion[];
}

export interface QuizShortResponse {
  id: number;
  title: string;
  summary: string;
  total_score: number;
  owner: {
    id: number;
    short_name: string;
    role: UserRole;
  };
  subject: {
    id: number;
    name: string;
  };
  group: Array<{
    id: number;
    name: string;
    curator: {
      id: number;
      short_name: string;
      role: UserRole;
    };
  }>;
}

export interface QuizFullResponse {
  id: number;
  title: string;
  summary: string;
  total_score: number;
  owner: {
    id: number;
    short_name: string;
    role: UserRole;
  };
  subject: {
    id: number;
    name: string;
  };
  group: Array<{
    id: number;
    name: string;
    curator: {
      id: number;
      short_name: string;
      role: UserRole;
    };
  }>;
  questions: QuizQuestion[];
}

export interface QuizProgressResponse {
  id: number;
  quiz: {
    id: number;
    title: string;
    summary: string;
    owner: UserShort;
    subject: SubjectResponse;
  };
  user: UserShort;
  status: ProgressStatus;
  score: number;
  start_date?: string;
  completed_date?: string;
}

export interface AnswerPatchRequest {
  text: string;
}

export interface ApiError {
  error: string;
}

// Типы для попыток прохождения квизов (старый формат, оставляем для совместимости)
export interface QuizAttemptAnswer {
  id: string;
  question_id: string;
  answer: string | null;
  score: number;
  is_correct: boolean;
}

export interface QuizAttempt {
  id: string;
  user: {
    id: string;
    first_name: string;
    last_name: string;
    role: string;
  };
  answers: QuizAttemptAnswer[];
  score: number;
  started_at: string;
  ended_at: string;
}

export interface QuizAttemptResponse {
  attempt: QuizAttempt;
  quiz: {
    id: string;
    title: string;
    summary: string;
    owner: {
      id: string;
      first_name: string;
      last_name: string;
      role: string;
    };
    subject: {
      id: string;
      name: string;
    };
    total_score: number;
    deadline: string;
    max_attempts: number;
    created_at: string;
    content: Array<{
      id: string;
      text: string;
      score: number;
      type: string;
      details: any;
    }>;
  };
}

// Новые типы для статистики квиза через /quizzes/{quiz_id}/users
export interface AttemptItem {
  id: string; // UUID
  score: number;
  started_at: string; // ISO date string
  ended_at: string | null; // ISO date string или null
}

export interface UserLastAttempt {
  id: string; // UUID
  first_name: string;
  last_name: string;
  role: string;
  last_attempt: AttemptItem | null; // Может быть null если попыток нет
}

// Типы для прохождения квиза
export interface QuizQuestionContent {
  id: string; // UUID
  text: string;
  score: number;
  type: 'single' | 'multiple' | 'numeric';
  details: QuizQuestionDetails;
}

export interface StartAttemptResponse {
  attempt: {
    id: string; // UUID
    started_at: string; // ISO date string
  };
  quiz: {
    id: string; // UUID
    title: string;
    summary: string;
    owner: {
      id: string;
      first_name: string;
      last_name: string;
      role: string;
    };
    subject: {
      id: string;
      name: string;
    };
    total_score: number;
    deadline: string; // ISO date string
    max_attempts: number;
    created_at: string; // ISO date string
    content: QuizQuestionContent[];
  };
}

/** ID попытки из ответа POST /attempts (только id). */
function extractAttemptIdFromCreateResponse(raw: unknown): string | null {
  if (!raw || typeof raw !== 'object') return null;
  const readId = (obj: Record<string, unknown>): string | null => {
    const v = obj.id ?? obj.attempt_id;
    if (typeof v === 'string' && v.trim() !== '') return v.trim();
    if (typeof v === 'number' && Number.isFinite(v)) return String(v);
    return null;
  };
  const o = raw as Record<string, unknown>;
  const direct = readId(o);
  if (direct) return direct;
  const att = o.attempt;
  if (att && typeof att === 'object') {
    const inner = readId(att as Record<string, unknown>);
    if (inner) return inner;
  }
  const nested = o.data;
  if (nested && typeof nested === 'object') {
    const d = nested as Record<string, unknown>;
    const fromData = readId(d);
    if (fromData) return fromData;
    const da = d.attempt;
    if (da && typeof da === 'object') {
      return readId(da as Record<string, unknown>);
    }
  }
  return null;
}

function withNormalizedQuizContent(res: StartAttemptResponse): StartAttemptResponse {
  const list = res.quiz?.content;
  if (!Array.isArray(list) || list.length === 0) return res;
  const content = list.map((item, i) => normalizeContentQuestion(item as unknown, i));
  return { ...res, quiz: { ...res.quiz, content } };
}

function normalizeStartAttemptResponse(raw: unknown): StartAttemptResponse {
  if (!raw || typeof raw !== 'object') {
    throw new Error('Пустой ответ при загрузке попытки');
  }
  const o = raw as Record<string, unknown>;
  if ('attempt' in o && 'quiz' in o) {
    return withNormalizedQuizContent(raw as StartAttemptResponse);
  }
  const inner = o.data;
  if (inner && typeof inner === 'object') {
    const d = inner as Record<string, unknown>;
    if ('attempt' in d && 'quiz' in d) {
      return withNormalizedQuizContent(inner as StartAttemptResponse);
    }
  }
  throw new Error('Некорректный ответ: ожидаются поля attempt и quiz');
}

function unwrapAttemptObject(raw: unknown): Record<string, unknown> | null {
  if (!raw || typeof raw !== 'object') return null;
  const o = raw as Record<string, unknown>;
  const a = o.attempt;
  if (a && typeof a === 'object') return a as Record<string, unknown>;
  const data = o.data;
  if (data && typeof data === 'object') {
    const inner = (data as Record<string, unknown>).attempt;
    if (inner && typeof inner === 'object') return inner as Record<string, unknown>;
  }
  return null;
}

function readQuizIdFromAttempt(attempt: Record<string, unknown>): string | null {
  return pickString(
    attempt.quiz_id as string,
    attempt.quizId as string
  );
}

function asStringArray(v: unknown): string[] {
  if (!Array.isArray(v)) return [];
  return v.map((x) => String(x));
}

/**
 * Приводит вопрос к формату UI: бэкенд отдаёт title, options и correct на верхнем уровне,
 * либо вложенный details (старый формат).
 */
function normalizeContentQuestion(item: unknown, index: number): QuizQuestionContent {
  if (!item || typeof item !== 'object') {
    throw new Error(`Некорректный вопрос #${index + 1} в квизе`);
  }
  const it = item as Record<string, unknown>;
  const id = String(it.id ?? `q-${index}`);
  const text = String(it.text ?? it.title ?? '');
  const score = Number(it.score ?? 0);

  const rawType = it.type;
  const correctTop = it.correct;
  const optionsRaw = it.options;
  const nested =
    it.details && typeof it.details === 'object' ? (it.details as Record<string, unknown>) : null;

  const hasOptions = Array.isArray(optionsRaw) && optionsRaw.length > 0;
  const hasNestedOptions = nested && Array.isArray(nested.options) && (nested.options as unknown[]).length > 0;

  let type: QuizQuestionContent['type'];
  if (rawType === 'multiple' || rawType === 'numeric' || rawType === 'single') {
    type = rawType;
  } else if (typeof correctTop === 'number' && !hasOptions && !hasNestedOptions) {
    type = 'numeric';
  } else if (Array.isArray(correctTop)) {
    type = 'multiple';
  } else if (hasOptions || hasNestedOptions) {
    type = 'single';
  } else {
    type = 'numeric';
  }

  let details: QuizQuestionDetails;

  if (type === 'numeric') {
    const c =
      typeof correctTop === 'number'
        ? correctTop
        : nested && typeof nested.correct === 'number'
          ? nested.correct
          : Number(correctTop ?? nested?.correct ?? 0);
    details = { correct: Number.isFinite(c) ? c : 0 };
  } else if (type === 'multiple') {
    const options = asStringArray(optionsRaw ?? nested?.options);
    const correctArr = Array.isArray(correctTop)
      ? asStringArray(correctTop)
      : asStringArray(nested?.correct);
    details = { options, correct: correctArr };
  } else {
    const options = asStringArray(optionsRaw ?? nested?.options);
    let correctStr = '';
    if (typeof correctTop === 'string') {
      correctStr = correctTop;
    } else if (Array.isArray(correctTop) && correctTop.length > 0) {
      correctStr = String(correctTop[0]);
    } else if (nested && typeof nested.correct === 'string') {
      correctStr = nested.correct;
    }
    details = { options, correct: correctStr };
  }

  return { id, text, score, type, details };
}

/** Догружает квиз для прохождения: GET /quizzes/:id (UUID). */
async function fetchQuizForTake(quizId: string): Promise<StartAttemptResponse['quiz']> {
  const qid = quizId.trim();
  if (!qid) throw new Error('Пустой quiz_id для загрузки квиза');

  const raw = await apiRequest<unknown>(`/quizzes/${encodeURIComponent(qid)}`, {
    method: 'GET',
  });

  let quiz: unknown = raw;
  if (raw && typeof raw === 'object' && 'quiz' in raw) {
    quiz = (raw as { quiz: unknown }).quiz;
  }

  if (!quiz || typeof quiz !== 'object') {
    throw new Error('Пустой ответ GET /quizzes/:id');
  }

  const o = quiz as Record<string, unknown>;
  const id = String(o.id ?? qid);
  const title = String(o.title ?? '');
  const summary = String(o.summary ?? '');
  const total_score = Number(o.total_score ?? 0);
  const deadline = String(o.deadline ?? '');
  const max_attempts = Number(o.max_attempts ?? 1);
  const created_at = String(o.created_at ?? deadline);

  const own = o.owner;
  let owner: StartAttemptResponse['quiz']['owner'] = {
    id: '',
    first_name: '',
    last_name: '',
    role: 'teacher',
  };
  if (own && typeof own === 'object') {
    const ow = own as Record<string, unknown>;
    if ('first_name' in ow || 'last_name' in ow) {
      owner = {
        id: String(ow.id ?? ''),
        first_name: String(ow.first_name ?? ''),
        last_name: String(ow.last_name ?? ''),
        role: String(ow.role ?? 'teacher'),
      };
    } else if ('short_name' in ow) {
      owner = {
        id: String(ow.id ?? ''),
        first_name: String(ow.short_name ?? ''),
        last_name: '',
        role: String(ow.role ?? 'teacher'),
      };
    }
  }

  const sub = o.subject;
  const subject: StartAttemptResponse['quiz']['subject'] = {
    id: sub && typeof sub === 'object' ? String((sub as Record<string, unknown>).id ?? '') : '',
    name: sub && typeof sub === 'object' ? String((sub as Record<string, unknown>).name ?? '') : '',
  };

  let content: QuizQuestionContent[] = [];
  const contentRaw = o.content;
  if (Array.isArray(contentRaw) && contentRaw.length > 0) {
    content = contentRaw.map((item, i) => normalizeContentQuestion(item, i));
  } else {
    const questions = o.questions;
    if (Array.isArray(questions) && questions.length > 0) {
      content = questions.map((item, i) => normalizeContentQuestion(item, i));
    }
  }

  return {
    id,
    title,
    summary,
    owner,
    subject,
    total_score,
    deadline,
    max_attempts,
    created_at,
    content,
  };
}

/** Собирает StartAttemptResponse, если бэкенд отдал только attempt (+ quiz_id). */
async function mergeAttemptWithFetchedQuiz(raw: unknown, fallbackAttemptId?: string): Promise<StartAttemptResponse> {
  const attemptObj = unwrapAttemptObject(raw);
  if (!attemptObj) {
    throw new Error('Некорректный ответ: нет объекта attempt');
  }
  const quizId = readQuizIdFromAttempt(attemptObj);
  if (!quizId) {
    throw new Error('В ответе нет quiz и у attempt нет quiz_id — нечего подгрузить');
  }
  const quizPart = await fetchQuizForTake(quizId);
  return withNormalizedQuizContent({
    attempt: {
      id: String(attemptObj.id ?? fallbackAttemptId ?? ''),
      started_at: String(attemptObj.started_at ?? ''),
    },
    quiz: quizPart,
  });
}

export interface FinishAttemptAnswer {
  id: string; // UUID
  question_id: string; // UUID
  answer: string | string[] | number; // Зависит от типа вопроса
  score: number;
  is_correct: boolean;
}

export interface FinishAttemptResponse {
  attempt: {
    id: string; // UUID
    user: {
      id: string;
      first_name: string;
      last_name: string;
      role: string;
    };
    answers: FinishAttemptAnswer[];
    /** Итоговые набранные баллы (бэкенд может назвать total_score) */
    score: number;
    started_at: string; // ISO date string
    ended_at: string; // ISO date string
  };
  quiz: {
    id: string; // UUID
    title: string;
    summary: string;
    owner: {
      id: string;
      first_name: string;
      last_name: string;
      role: string;
    };
    subject: {
      id: string;
      name: string;
    };
    total_score: number;
    deadline: string;
    max_attempts: number;
    created_at: string;
    content: QuizQuestionContent[];
  };
}

function normalizeFinishAnswerItem(item: unknown, index: number): FinishAttemptAnswer {
  if (!item || typeof item !== 'object') {
    return {
      id: `answer-${index}`,
      question_id: '',
      answer: '',
      score: 0,
      is_correct: false,
    };
  }
  const x = item as Record<string, unknown>;
  const ans = x.answer;
  const answerVal: string | string[] | number =
    ans === null || ans === undefined ? '' : (ans as string | string[] | number);
  const ic = x.is_correct;
  const rawScore = Number(x.score ?? 0);
  const isCorrect =
    ic === true || ic === 'true' || ic === 1 || rawScore > 0;
  return {
    id: String(x.id ?? `answer-${index}`),
    question_id: String(x.question_id ?? ''),
    answer: answerVal,
    score: rawScore,
    is_correct: isCorrect,
  };
}

function extractAnswersArray(att: Record<string, unknown>): unknown[] | null {
  let raw: unknown = att.answers ?? att.Answers;
  if (typeof raw === 'string') {
    try {
      raw = JSON.parse(raw) as unknown;
    } catch {
      return null;
    }
  }
  if (!Array.isArray(raw)) return null;
  return raw;
}

function answersMatchForGrading(
  q: QuizQuestionContent,
  user: string | string[] | number
): boolean {
  if (q.type === 'numeric') {
    const c = (q.details as QuizQuestionDetailsNumeric).correct;
    const u = typeof user === 'number' ? user : parseFloat(String(user).trim().replace(',', '.'));
    if (Number.isNaN(u)) return false;
    return Number(u) === Number(c);
  }
  if (q.type === 'single') {
    const c = (q.details as QuizQuestionDetailsSingle).correct;
    return String(user).trim() === String(c).trim();
  }
  const mul = q.details as QuizQuestionDetailsMultiple;
  const correct = [...mul.correct].map(String).sort();
  const uArr = (Array.isArray(user) ? user : [user]).map(String).sort();
  if (uArr.length !== correct.length) return false;
  return uArr.every((v, idx) => v === correct[idx]);
}

function buildGradedClientAnswers(
  clientAnswers: Array<{ question_id: string; answer: string | string[] | number }>,
  quiz: StartAttemptResponse['quiz']
): FinishAttemptAnswer[] {
  return clientAnswers.map((ca, i) => {
    const q = quiz.content.find((c) => c.id === ca.question_id);
    const ok = q ? answersMatchForGrading(q, ca.answer) : false;
    const pts = q && ok ? q.score : 0;
    return {
      id: `client-${i}`,
      question_id: ca.question_id,
      answer: ca.answer,
      score: pts,
      is_correct: ok,
    };
  });
}

/** Ответ GET /attempts/:id с непустым answers (после finish). */
function getAttemptWithAnswersFromPayload(raw: unknown): Record<string, unknown> | null {
  const inner = unwrapAttemptObject(raw);
  if (inner) {
    const arr = extractAnswersArray(inner);
    if (arr && arr.length > 0) return inner;
  }
  if (raw && typeof raw === 'object') {
    const o = raw as Record<string, unknown>;
    if (o.id != null || o.quiz_id != null) {
      const arr = extractAnswersArray(o);
      if (arr && arr.length > 0) return o;
    }
  }
  return null;
}

async function fetchAttemptDetailAfterFinish(attemptId: string): Promise<unknown | null> {
  const id = String(attemptId).trim();
  if (!id) return null;
  let last: unknown | null = null;
  const maxTries = 6;
  for (let i = 0; i < maxTries; i++) {
    if (i > 0) {
      await new Promise((r) => setTimeout(r, 180 * i));
    }
    try {
      last = await apiRequest<unknown>(`/attempts/${encodeURIComponent(id)}`, {
        method: 'GET',
      });
      if (getAttemptWithAnswersFromPayload(last)) {
        return last;
      }
    } catch {
      last = null;
    }
  }
  return last;
}

/** POST finish часто без answers; подменяем телом GET /attempts/:id. */
function preferAttemptDetailForFinish(finishRaw: unknown, getRaw: unknown | null): unknown {
  const det = getAttemptWithAnswersFromPayload(getRaw);
  if (!det) return finishRaw;
  if (getRaw && typeof getRaw === 'object' && 'attempt' in getRaw) {
    return getRaw;
  }
  return { attempt: det };
}

/**
 * Бэкенд часто отдаёт только { attempt: { ..., total_score } } без quiz и answers —
 * подставляем квиз с клиента и приводим score/answers.
 */
function normalizeFinishAttemptResponse(
  raw: unknown,
  ctx: {
    fallbackQuiz?: StartAttemptResponse['quiz'];
    clientAnswers?: Array<{ question_id: string; answer: string | string[] | number }>;
  }
): FinishAttemptResponse {
  if (!raw || typeof raw !== 'object') {
    throw new Error('Пустой ответ POST .../finish');
  }
  let root = raw as Record<string, unknown>;
  if (root.data && typeof root.data === 'object') {
    root = root.data as Record<string, unknown>;
  }

  if (root.attempt && root.quiz && typeof root.quiz === 'object') {
    const ar = root.attempt as Record<string, unknown>;
    const scoreFull =
      typeof ar.score === 'number'
        ? ar.score
        : typeof ar.total_score === 'number'
          ? ar.total_score
          : 0;
    const userRawF = ar.user;
    const userF =
      userRawF && typeof userRawF === 'object'
        ? {
            id: String((userRawF as Record<string, unknown>).id ?? ''),
            first_name: String((userRawF as Record<string, unknown>).first_name ?? ''),
            last_name: String((userRawF as Record<string, unknown>).last_name ?? ''),
            role: String((userRawF as Record<string, unknown>).role ?? 'student'),
          }
        : {
            id: String(ar.user_id ?? ''),
            first_name: '',
            last_name: '',
            role: 'student',
          };
    let answersF: FinishAttemptAnswer[] = [];
    const ansArrF = extractAnswersArray(ar);
    if (ansArrF && ansArrF.length > 0) {
      answersF = ansArrF.map((item, i) => normalizeFinishAnswerItem(item, i));
    } else if (ctx.clientAnswers?.length) {
      const qGrade = withNormalizedQuizContent({
        attempt: {
          id: String(ar.id ?? ''),
          started_at: String(ar.started_at ?? ''),
        },
        quiz: root.quiz as StartAttemptResponse['quiz'],
      }).quiz;
      answersF = buildGradedClientAnswers(ctx.clientAnswers, qGrade);
    }
    const qz = withNormalizedQuizContent({
      attempt: {
        id: String(ar.id ?? ''),
        started_at: String(ar.started_at ?? ''),
      },
      quiz: root.quiz as StartAttemptResponse['quiz'],
    });
    return {
      attempt: {
        id: String(ar.id ?? ''),
        user: userF,
        answers: answersF,
        score: scoreFull,
        started_at: String(ar.started_at ?? ''),
        ended_at: String(ar.ended_at ?? ''),
      },
      quiz: qz.quiz,
    };
  }

  const att = root.attempt;
  if (!att || typeof att !== 'object') {
    throw new Error('В ответе finish нет объекта attempt');
  }
  const a = att as Record<string, unknown>;

  const id = String(a.id ?? '');
  const started_at = String(a.started_at ?? '');
  const ended_at = String(a.ended_at ?? '');
  const score =
    typeof a.score === 'number'
      ? a.score
      : typeof a.total_score === 'number'
        ? a.total_score
        : Number(a.score ?? a.total_score ?? 0);

  const userRaw = a.user;
  const user =
    userRaw && typeof userRaw === 'object'
      ? {
          id: String((userRaw as Record<string, unknown>).id ?? ''),
          first_name: String((userRaw as Record<string, unknown>).first_name ?? ''),
          last_name: String((userRaw as Record<string, unknown>).last_name ?? ''),
          role: String((userRaw as Record<string, unknown>).role ?? 'student'),
        }
      : {
          id: String(a.user_id ?? ''),
          first_name: '',
          last_name: '',
          role: 'student',
        };

  let answers: FinishAttemptAnswer[] = [];
  const ansArr = extractAnswersArray(a);
  if (ansArr && ansArr.length > 0) {
    answers = ansArr.map((item, i) => normalizeFinishAnswerItem(item, i));
  } else if (ctx.clientAnswers?.length && ctx.fallbackQuiz) {
    answers = buildGradedClientAnswers(ctx.clientAnswers, ctx.fallbackQuiz);
  } else if (ctx.clientAnswers?.length) {
    answers = ctx.clientAnswers.map((ca, i) => ({
      id: `client-${i}`,
      question_id: ca.question_id,
      answer: ca.answer,
      score: 0,
      is_correct: false,
    }));
  }

  if (!ctx.fallbackQuiz) {
    throw new Error('Ответ finish без quiz: передайте текущий квиз с клиента');
  }

  const quiz = withNormalizedQuizContent({
    attempt: { id, started_at },
    quiz: ctx.fallbackQuiz,
  }).quiz;

  return {
    attempt: {
      id,
      user,
      answers,
      score,
      started_at,
      ended_at,
    },
    quiz,
  };
}

// Утилиты для работы с авторизацией (user_id и role)
export interface AuthData {
  user_id: string;
  role: string;
}

export const getAccessToken = (): string | null => localStorage.getItem('access_token');

export const getAuthData = (): AuthData | null => {
  const userId = localStorage.getItem('user_id');
  const role = localStorage.getItem('user_role');
  if (userId && role) {
    return { user_id: userId, role };
  }
  const token = getAccessToken();
  if (token) {
    const payload = readJwtPayload(token);
    const uid = payload?.sub ?? payload?.user_id;
    if (uid) {
      return {
        user_id: String(uid),
        role: String(payload?.role ?? 'student'),
      };
    }
  }
  return null;
};

export const setAuthData = (userId: string, role: string): void => {
  localStorage.setItem('user_id', userId);
  localStorage.setItem('user_role', role);
};

export const removeAuthData = (): void => {
  localStorage.removeItem('user_id');
  localStorage.removeItem('user_role');
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
};

export const isAuthenticated = (): boolean => {
  return !!getAccessToken();
};

// Обратная совместимость (для старых методов)
export const getAuthToken = (): string | null => getAccessToken();

export const setAuthToken = (token: string): void => {
  localStorage.setItem('access_token', token.trim());
};

export const removeAuthToken = (): void => {
  removeAuthData();
};

export const getRefreshToken = (): string | null => {
  return localStorage.getItem('refresh_token');
};

export const setRefreshToken = (token: string): void => {
  localStorage.setItem('refresh_token', token);
};

export const removeRefreshToken = (): void => {
  localStorage.removeItem('refresh_token');
};

let refreshInFlight: Promise<boolean> | null = null;

async function tryRefreshAccessToken(): Promise<boolean> {
  const current = localStorage.getItem('refresh_token');
  if (!current) return false;
  if (refreshInFlight) return refreshInFlight;
  const base = API_BASE_URL;
  if (!base) return false;

  refreshInFlight = (async () => {
    try {
      const res = await fetch(`${base}/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: current }),
      });
      if (!res.ok) return false;
      const data: unknown = await res.json();
      const access = extractAccessTokenFromPayload(data);
      const nextRefresh = extractRefreshTokenFromPayload(data);
      if (access) localStorage.setItem('access_token', access);
      if (nextRefresh) localStorage.setItem('refresh_token', nextRefresh);
      return !!access;
    } catch {
      return false;
    } finally {
      refreshInFlight = null;
    }
  })();

  return refreshInFlight;
}

// Базовая функция для запросов
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {},
  _retryAfterRefresh = false
): Promise<T> {
  // Локальный режим разработки без бэкенда
  if (LOCAL_DEV_MODE) {
    console.warn(`[API] Локальный режим разработки: запрос ${options.method || 'GET'} ${endpoint} пропущен`);
    // Возвращаем пустые данные в зависимости от типа запроса
    if (endpoint.includes('/auth/login')) {
      return {
        id: 'mock_user_' + Date.now(),
        role: 'student',
        access_token: 'dev-mock-access',
        refresh_token: 'dev-mock-refresh',
      } as T;
    }
    if (endpoint.includes('/users/me')) {
      // Мок-данные для текущего пользователя
      return {
        id: 1,
        first_name: 'Тестовый',
        last_name: 'Пользователь',
        role: 'student',
      } as T;
    }
    if (endpoint.includes('/users') && options.method !== 'POST' && !endpoint.includes('/users/')) {
      return { users: [] } as T;
    }
    if (endpoint.includes('/groups') && !endpoint.includes('/groups/')) {
      return { groups: [] } as T;
    }
    if (endpoint.includes('/subjects')) {
      return { subjects: [] } as T;
    }
    if (endpoint.includes('/quizzes') && !endpoint.includes('/quizzes/')) {
      return { quizzes: [] } as T;
    }
    const attemptsDetail = /^\/attempts\/([^/]+)$/.exec(endpoint);
    if (attemptsDetail && (options.method === 'GET' || options.method === undefined)) {
      const aid = attemptsDetail[1];
      return {
        attempt: { id: aid, started_at: new Date().toISOString() },
        quiz: {
          id: 'mock-quiz',
          title: 'Мок-квиз',
          summary: 'Локальный режим',
          owner: { id: '1', first_name: 'Препод', last_name: '', role: 'teacher' },
          subject: { id: '1', name: 'Предмет' },
          total_score: 10,
          deadline: new Date().toISOString(),
          max_attempts: 3,
          created_at: new Date().toISOString(),
          content: [],
        },
      } as T;
    }
    if (endpoint === '/attempts' && options.method === 'POST') {
      return { id: 'mock-attempt-id' } as T;
    }
    const finishMatch = /^\/attempts\/([^/]+)\/finish$/.exec(endpoint);
    if (finishMatch && options.method === 'POST') {
      const aid = finishMatch[1];
      return {
        attempt: {
          id: aid,
          quiz_id: 'mock-quiz',
          started_at: new Date().toISOString(),
          ended_at: new Date().toISOString(),
          total_score: 0,
          user_id: 'mock-user',
        },
      } as T;
    }
    if (endpoint.includes('/progress')) {
      return { progress: [] } as T;
    }
    if (endpoint.includes('/groups/') && endpoint.includes('/students')) {
      return { students: [] } as T;
    }
    if (endpoint.includes('/quizzes/') && endpoint.includes('/attempts')) {
      return { attempts: [] } as T;
    }
    if (endpoint.includes('/users/') && endpoint.includes('/quizzes')) {
      return { quizzes: [] } as T;
    }
    // Для POST/PATCH/DELETE просто возвращаем пустой объект
    return {} as T;
  }

  if (!API_BASE_URL) {
    throw new Error('API_BASE_URL не настроен. Установите VITE_API_URL или включите proxy в vite.config.ts');
  }

  const authData = getAuthData();
  const method = options.method || 'GET';
  const bearer = localStorage.getItem('access_token');
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> || {}),
  };

  if (bearer) {
    headers['Authorization'] = `Bearer ${bearer}`;
  }

  const skipLegacyUserAuth = isPublicAuthPath(endpoint);
  const useLegacyUserId = authData && !bearer && !skipLegacyUserAuth;

  // Без JWT: для GET отправляем user_id и role в заголовках (контракт старого бэкенда)
  if (useLegacyUserId && method === 'GET') {
    headers['X-User-Id'] = authData.user_id;
    headers['X-User-Role'] = authData.role;
  }

  try {
    // Подготавливаем тело запроса с user_id и role для POST/PATCH/PUT запросов
    let requestBody = options.body;
    
    // С JWT не подмешиваем user_id в тело — идентификация через Authorization.
    // Публичные /auth/* не трогаем (иначе в /auth/login попадали старые user_id и ломали бэкенд).
    if (
      useLegacyUserId &&
      (method === 'POST' || method === 'PATCH' || method === 'PUT')
    ) {
      try {
        const bodyData = requestBody ? JSON.parse(requestBody as string) : {};
        bodyData.user_id = authData.user_id;
        // Добавляем role только если его нет в теле запроса (чтобы не перезаписывать роль при создании пользователя)
        if (!bodyData.role) {
          bodyData.role = authData.role;
        }
        requestBody = JSON.stringify(bodyData);
      } catch (e) {
        // Если не удалось распарсить тело, создаем новое
        requestBody = JSON.stringify({
          user_id: authData.user_id,
          role: authData.role,
        });
      }
    }
    
    // Отладочная информация
    if (import.meta.env.DEV) {
      const bodyForLog = requestBody ? (() => {
        try {
          return JSON.parse(requestBody as string);
        } catch {
          return requestBody;
        }
      })() : undefined;
      console.log(`[API] ${method} ${endpoint}`, {
        hasAuth: !!authData,
        hasBearer: !!bearer,
        userId: authData?.user_id,
        role: authData?.role,
        url: `${API_BASE_URL}${endpoint}`,
        body: bodyForLog,
      });
    }
    
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
      body: requestBody,
    });

    // Обработка сетевых ошибок
    if (!response.ok) {
      // Если 401 - неавторизован
      if (response.status === 401) {
        if (
          !_retryAfterRefresh &&
          !isPublicAuthPath(endpoint) &&
          localStorage.getItem('refresh_token')
        ) {
          const refreshed = await tryRefreshAccessToken();
          if (refreshed) {
            return apiRequest<T>(endpoint, options, true);
          }
        }
        logout();
        throw new Error('Сессия истекла. Пожалуйста, войдите снова.');
      }

      let errorMessage = `HTTP error! status: ${response.status}`;
      try {
        const errorData: ApiError = await response.json();
        errorMessage = errorData.error || errorMessage;
      } catch {
        // Если не удалось распарсить ошибку, используем стандартное сообщение
        if (response.status === 404) {
          errorMessage = 'Ресурс не найден';
        } else if (response.status === 500) {
          errorMessage = 'Ошибка сервера';
        }
      }
      throw new Error(errorMessage);
    }

    // Если ответ пустой (например, для DELETE запросов)
    if (response.status === 204 || response.headers.get('content-length') === '0') {
      return {} as T;
    }

    return await response.json();
  } catch (error) {
    // Обработка сетевых ошибок (CORS, таймаут, и т.д.)
    if (error instanceof TypeError && error.message === 'Failed to fetch') {
      const backendUrl = import.meta.env.DEV
        ? `${import.meta.env.VITE_API_URL ?? 'см. VITE_API_URL в frontend/.env'} (или прокси /api)`
        : API_BASE_URL;
      throw new Error(`Не удалось подключиться к серверу. Проверьте, что бэкенд доступен на ${backendUrl}`);
    }
    throw error;
  }
}

// API методы

// Авторизация
export const login = async (credentials: LoginRequest): Promise<{ user_id: string; role: string }> => {
  // Локальный режим разработки - мок-авторизация
  if (LOCAL_DEV_MODE) {
    console.warn('[API] Локальный режим: используется мок-авторизация');
    const mockUserId = 'mock_user_' + Date.now();
    const mockRole = 'student';
    setAuthData(mockUserId, mockRole);
    localStorage.setItem('access_token', 'dev-mock-bearer');
    return { user_id: mockUserId, role: mockRole };
  }

  removeAuthData();

  const response = await apiRequest<LoginApiResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(credentials),
  });

  const accessToken = extractAccessTokenFromPayload(response);
  const refreshTok = extractRefreshTokenFromPayload(response);
  if (accessToken) localStorage.setItem('access_token', accessToken);
  if (refreshTok) localStorage.setItem('refresh_token', refreshTok);

  let userId = [response.id, response.user_id, response.user?.id]
    .filter((v) => v != null && String(v).trim() !== '')
    .map((v) => String(v))[0] ?? '';

  let role = String(response.role ?? response.user?.role ?? '').trim() || 'student';

  if (!userId && accessToken) {
    const payload = readJwtPayload(accessToken);
    if (payload) {
      userId = String(payload.sub ?? payload.user_id ?? '').trim();
      if (payload.role) role = String(payload.role);
    }
  }

  if (!accessToken) {
    throw new Error(
      'Сервер не вернул токен авторизации. Ожидаются поля access_token, token или объект jwt в ответе /auth/login'
    );
  }

  if (!userId) {
    throw new Error('Некорректный ответ при входе: нет id пользователя в теле ответа и в JWT');
  }

  setAuthData(userId, role);

  return { user_id: userId, role };
};

// Пользователи
export const createUser = async (userData: UserCreateRequest): Promise<void> => {
  await apiRequest<void>('/users', {
    method: 'POST',
    body: JSON.stringify(userData),
  });
};

export const getUsers = async (): Promise<UserShort[]> => {
  const response = await apiRequest<{ users: UserShort[] }>('/users');
  return response.users || [];
};

export const getCurrentUser = async (): Promise<UserFull> => {
  if (!getAccessToken()) {
    throw new Error('Пользователь не авторизован');
  }
  const authData = getAuthData();
  if (!authData) {
    throw new Error('Пользователь не авторизован');
  }
  const response = await apiRequest<{ user: UserFull }>(`/users/${authData.user_id}`, {
    method: 'GET',
  });
  // Проверяем структуру ответа
  if (response && 'user' in response) {
    return response.user;
  }
  // Если ответ пришел напрямую как UserFull
  return response as unknown as UserFull;
};

export const getUserById = async (userId: number): Promise<UserFull> => {
  const response = await apiRequest<{ user: UserFull }>(`/users/${userId}`);
  // Проверяем структуру ответа
  if (response && 'user' in response) {
    return response.user;
  }
  // Если ответ пришел напрямую как UserFull
  return response as unknown as UserFull;
};

// Группы
export const createGroup = async (groupData: GroupCreateRequest): Promise<void> => {
  await apiRequest<void>('/groups', {
    method: 'POST',
    body: JSON.stringify(groupData),
  });
};

export const getGroups = async (): Promise<GroupResponse[]> => {
  const response = await apiRequest<{ groups: GroupResponse[] }>('/groups');
  return response.groups || [];
};

export const addStudentsToGroup = async (
  groupId: string | number,
  studentIds: GroupAddStudentsRequest
): Promise<void> => {
  const groupIdStr = String(groupId);
  await apiRequest<void>(`/groups/${groupIdStr}/students`, {
    method: 'POST',
    body: JSON.stringify(studentIds),
  });
};

// Получить студентов группы
export const getGroupStudents = async (groupId: string | number): Promise<UserShort[]> => {
  const groupIdStr = String(groupId);
  const response = await apiRequest<{ students: UserShort[] }>(`/groups/${groupIdStr}/students`);
  return response.students || [];
};

// Предметы
export const createSubject = async (subjectData: SubjectCreateRequest): Promise<void> => {
  await apiRequest<void>('/subjects', {
    method: 'POST',
    body: JSON.stringify(subjectData),
  });
};

export const getSubjects = async (): Promise<SubjectResponse[]> => {
  const response = await apiRequest<{ subjects: SubjectResponse[] }>('/subjects');
  return response.subjects || [];
};

// Квизы
export const createQuiz = async (quizData: QuizCreateRequest): Promise<void> => {
  await apiRequest<void>('/quizzes', {
    method: 'POST',
    body: JSON.stringify(quizData),
  });
};

export const getQuizzes = async (): Promise<QuizShortResponse[]> => {
  const response = await apiRequest<{ quizzes: QuizShortResponse[] }>('/quizzes');
  return response.quizzes || [];
};

export const getQuizById = async (quizId: number): Promise<QuizFullResponse> => {
  const response = await apiRequest<{ quiz: QuizFullResponse }>(`/quizzes/${quizId}`);
  // Проверяем структуру ответа
  if (response && 'quiz' in response) {
    return response.quiz;
  }
  // Если ответ пришел напрямую как QuizFullResponse
  return response as unknown as QuizFullResponse;
};

export const deleteQuiz = async (quizId: number): Promise<void> => {
  await apiRequest<void>(`/quizzes/${quizId}`, {
    method: 'DELETE',
  });
};

// Получение попыток прохождения квиза (старый эндпоинт)
export const getQuizAttempts = async (quizId: string | number): Promise<QuizAttemptResponse[]> => {
  const response = await apiRequest<{ attempts: QuizAttemptResponse[] }>(`/quizzes/${quizId}/attempts`);
  return response.attempts || [];
};

// Получение студентов с их последними попытками для квиза
export const getQuizUsers = async (quizId: string | number): Promise<UserLastAttempt[]> => {
  const response = await apiRequest<UserLastAttempt[]>(`/quizzes/${quizId}/users`);
  return response || [];
};

/** GET /attempts/:id — попытка и квиз (квиз может догружаться GET /quizzes/:quiz_id) */
export const getAttemptById = async (attemptId: string): Promise<StartAttemptResponse> => {
  const id = String(attemptId).trim();
  if (!id) throw new Error('Не указан id попытки');
  const response = await apiRequest<unknown>(`/attempts/${encodeURIComponent(id)}`, {
    method: 'GET',
  });
  try {
    return normalizeStartAttemptResponse(response);
  } catch {
    return mergeAttemptWithFetchedQuiz(response, id);
  }
};

/** Создать попытку: POST /attempts { user_id, quiz_id } → id; затем GET /attempts/:id */
export const startQuizAttempt = async (quizId: string | number): Promise<StartAttemptResponse> => {
  const authData = getAuthData();
  if (!authData) {
    throw new Error('Пользователь не авторизован');
  }
  const quiz_id = String(quizId);
  const user_id = authData.user_id;
  const created = await apiRequest<unknown>('/attempts', {
    method: 'POST',
    body: JSON.stringify({ user_id, quiz_id }),
  });
  try {
    return normalizeStartAttemptResponse(created);
  } catch {
    /* нет полного { attempt, quiz } */
  }
  try {
    return await mergeAttemptWithFetchedQuiz(created);
  } catch {
    /* нет вложенного attempt с quiz_id — цепочка через id */
  }
  const attemptId = extractAttemptIdFromCreateResponse(created);
  if (!attemptId) {
    throw new Error('Сервер не вернул id попытки после POST /attempts');
  }
  return getAttemptById(attemptId);
};

// Завершить попытку прохождения квиза
// Ответы должны быть отправлены отдельно через PATCH /attempts/{attempt_id}/answers/{answer_id}
// или включены в тело запроса finish (зависит от реализации бэкенда)
export const finishQuizAttempt = async (
  attemptId: string,
  answers?: Array<{ question_id: string; answer: string | string[] | number }>,
  fallbackQuiz?: StartAttemptResponse['quiz']
): Promise<FinishAttemptResponse> => {
  const id = String(attemptId).trim();
  const body = answers ? { answers } : undefined;
  const finishRaw = await apiRequest<unknown>(`/attempts/${encodeURIComponent(id)}/finish`, {
    method: 'POST',
    body: body ? JSON.stringify(body) : undefined,
  });
  const detailRaw = await fetchAttemptDetailAfterFinish(id);
  const sourceRaw = preferAttemptDetailForFinish(finishRaw, detailRaw);
  return normalizeFinishAttemptResponse(sourceRaw, { fallbackQuiz, clientAnswers: answers });
};

// Прогресс
export const getProgress = async (): Promise<QuizProgressResponse[]> => {
  const response = await apiRequest<{ progress: QuizProgressResponse[] }>('/progress');
  return response.progress || [];
};

// Получить квизы пользователя
export const getUserQuizzes = async (userId?: string): Promise<QuizShortResponse[]> => {
  const authData = getAuthData();
  let targetUserId = userId;
  
  // Если userId не передан, берем из authData
  if (!targetUserId && authData) {
    targetUserId = authData.user_id;
  }
  
  // Проверяем, что user_id валидный (не 'true', не пустая строка)
  const userIdStr = String(targetUserId || '');
  if (!targetUserId || userIdStr === 'true' || userIdStr === 'false' || userIdStr.trim() === '') {
    throw new Error('Пользователь не авторизован или user_id неверный');
  }
  
  // Убеждаемся, что user_id - это строка
  const userIdString = String(targetUserId);
  const response = await apiRequest<{ quizzes: QuizShortResponse[] }>(`/users/${userIdString}/quizzes`);
  return response.quizzes || [];
};

export const startProgress = async (progressId: number): Promise<void> => {
  await apiRequest<void>(`/progress/${progressId}/start`, {
    method: 'POST',
  });
};

export const finishProgress = async (progressId: number): Promise<void> => {
  await apiRequest<void>(`/progress/${progressId}/finish`, {
    method: 'POST',
  });
};

export const patchAnswer = async (
  progressId: number,
  answerId: number,
  answerData: AnswerPatchRequest
): Promise<void> => {
  await apiRequest<void>(`/progress/${progressId}/answer/${answerId}`, {
    method: 'PATCH',
    body: JSON.stringify(answerData),
  });
};

export const markAnswerCorrect = async (
  progressId: number,
  answerId: number
): Promise<void> => {
  await apiRequest<void>(`/progress/${progressId}/answer/${answerId}/correct`, {
    method: 'POST',
  });
};

export const markAnswerIncorrect = async (
  progressId: number,
  answerId: number
): Promise<void> => {
  await apiRequest<void>(`/progress/${progressId}/answer/${answerId}/incorrect`, {
    method: 'POST',
  });
};

// Выход
export const logout = (): void => {
  removeAuthData();
  localStorage.removeItem('user');
};

