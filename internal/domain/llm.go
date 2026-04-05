package domain

const SystemPrompt = `Ты — парсер данных для аварийной службы. Извлеки информацию из текста и верни ТОЛЬКО валидный JSON.
ПРАВИЛА:
1. Верни строго один JSON-объект, без пояснений, без markdown
2. Все строковые поля в двойных кавычках, числа без кавычек
3. Если поле неизвестно — используй: "" для строк, 0 для чисел
ФОРМАТ ОТВЕТА:
{
  "address": "строка",
  "latitude": число,
  "longitude": число,
  "priority": "critical"|"high"|"normal",
  "issue_type": "power_outage"|"equipment_failure"|"fire"|"other",
  "description": "строка",
  "confidence": число_от_0_до_1
}
ЗНАЧЕНИЯ ПОЛЕЙ:
• priority: "critical"=угроза жизни, "high"=важное оборудование, "normal"=обычная
• issue_type: "power_outage"=нет света, "equipment_failure"=искрение/дым, "fire"=пожар, "other"=прочее
• latitude/longitude: заполняй ТОЛЬКО если в тексте есть явные координаты, иначе 0
• confidence: твоя уверенность (0.0–1.0)
`

type Request struct {
	Model             string    `json:"model"`
	Messages          []Message `json:"messages"`
	N                 int       `json:"n,omitempty"`
	Stream            bool      `json:"stream,omitempty"`
	MaxTokens         int       `json:"max_tokens,omitempty"`
	RepetitionPenalty float64   `json:"repetition_penalty,omitempty"`
	UpdateInterval    int       `json:"update_interval,omitempty"`
	Temperature       float64   `json:"temperature,omitempty"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
