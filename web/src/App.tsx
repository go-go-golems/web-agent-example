import { useCallback, useMemo, useState } from 'react';
import type { ComposerSlotProps, ChatWidgetRenderers } from '@pwchat/webchat';
import { ChatWidget } from '@pwchat/webchat';
import { ThinkingModeComposer } from './components/ThinkingModeComposer';
import { WebAgentThinkingModeCard } from './components/WebAgentThinkingModeCard';

export function App() {
  const [mode, setMode] = useState('fast');

  const buildOverrides = useCallback(() => {
    return {
      middlewares: [
        {
          name: 'webagent-thinking-mode',
          config: { mode },
        },
      ],
    };
  }, [mode]);

  const renderers: ChatWidgetRenderers = useMemo(
    () => ({
      webagent_thinking_mode: WebAgentThinkingModeCard,
    }),
    [],
  );

  const components = useMemo(
    () => ({
      Composer: (props: ComposerSlotProps) => (
        <ThinkingModeComposer {...props} mode={mode} onModeChange={setMode} />
      ),
    }),
    [mode],
  );

  return (
    <div className="app-shell">
      <div className="app-card">
        <ChatWidget renderers={renderers} components={components} buildOverrides={buildOverrides} />
      </div>
    </div>
  );
}
