import type { Meta, StoryObj } from '@storybook/react';
import { BlockCard } from './BlockCard';
import type { ParsedBlock } from '../types';

const meta: Meta<typeof BlockCard> = {
  title: 'Debug UI/BlockCard',
  component: BlockCard,
  parameters: {
    layout: 'padded',
  },
};

export default meta;
type Story = StoryObj<typeof BlockCard>;

const systemBlock: ParsedBlock = {
  index: 0,
  kind: 'system',
  role: 'system',
  payload: { text: 'You are a helpful assistant. Be concise and accurate.' },
  metadata: { 'geppetto.middleware@v1': 'system-prompt-mw' },
};

const userBlock: ParsedBlock = {
  index: 1,
  kind: 'user',
  role: 'user',
  payload: { text: 'What is the weather in Paris?' },
  metadata: {},
};

const toolCallBlock: ParsedBlock = {
  index: 2,
  id: 'tc_001',
  kind: 'tool_call',
  payload: {
    id: 'tc_001',
    name: 'get_weather',
    args: { location: 'Paris', units: 'celsius' },
  },
  metadata: {},
};

const toolUseBlock: ParsedBlock = {
  index: 3,
  kind: 'tool_use',
  payload: {
    id: 'tc_001',
    result: { temperature: 18, condition: 'cloudy', humidity: 65 },
  },
  metadata: {},
};

const llmTextBlock: ParsedBlock = {
  index: 4,
  kind: 'llm_text',
  role: 'assistant',
  payload: { text: 'The weather in Paris is currently 18°C and cloudy with 65% humidity.' },
  metadata: {},
};

const reasoningBlock: ParsedBlock = {
  index: 5,
  kind: 'reasoning',
  payload: { encrypted_content: '<encrypted>' },
  metadata: {},
};

export const System: Story = {
  args: { block: systemBlock },
};

export const User: Story = {
  args: { block: userBlock },
};

export const ToolCall: Story = {
  args: { block: toolCallBlock },
};

export const ToolUse: Story = {
  args: { block: toolUseBlock },
};

export const LLMText: Story = {
  args: { block: llmTextBlock },
};

export const Reasoning: Story = {
  args: { block: reasoningBlock },
};

export const NewBlock: Story = {
  args: { block: llmTextBlock, isNew: true },
};

export const Compact: Story = {
  args: { block: llmTextBlock, compact: true },
};

export const LongText: Story = {
  args: {
    block: {
      ...llmTextBlock,
      payload: {
        text: `The weather in Paris is currently 18°C and cloudy with 65% humidity. 
        
This is a longer response that demonstrates how the component handles multi-line text content. The weather pattern is expected to continue for the next few days, with temperatures ranging from 15°C to 20°C.

Key points:
1. Current temperature: 18°C
2. Condition: Cloudy
3. Humidity: 65%
4. Wind: Light breeze from the west

Would you like more detailed information about the forecast for the coming week?`,
      },
    },
  },
};

export const Expanded: Story = {
  args: {
    block: {
      ...systemBlock,
      metadata: {
        'geppetto.middleware@v1': 'system-prompt-mw',
        'geppetto.session_id@v1': 'sess_01234567890abcdef',
        'geppetto.inference_id@v1': 'inf_abcdef1234567890',
        'custom.config@v1': { mode: 'agent', tools_enabled: true, max_tokens: 4096 },
      },
    },
    expanded: true,
  },
};

export const WithRichMetadata: Story = {
  args: {
    block: {
      ...llmTextBlock,
      metadata: {
        'geppetto.middleware@v1': 'agent-mode-mw',
        'geppetto.usage@v1': { prompt_tokens: 1234, completion_tokens: 567, total_tokens: 1801 },
        'planning.context@v1': { run_id: 'plan_xyz', step: 3, status: 'executing' },
        'tool.calls@v1': ['get_weather', 'search_web'],
      },
    },
  },
};

export const AllBlockTypes: Story = {
  render: () => (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', maxWidth: '600px' }}>
      <BlockCard block={systemBlock} />
      <BlockCard block={userBlock} />
      <BlockCard block={toolCallBlock} />
      <BlockCard block={toolUseBlock} />
      <BlockCard block={llmTextBlock} />
      <BlockCard block={reasoningBlock} />
    </div>
  ),
};
