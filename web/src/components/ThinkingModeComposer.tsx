import { DefaultComposer, type ComposerSlotProps } from '@pwchat/webchat';

type Props = ComposerSlotProps & {
  mode: string;
  onModeChange: (next: string) => void;
  discoEnabled: boolean;
  onDiscoEnabledChange: (next: boolean) => void;
};

export function ThinkingModeComposer({ mode, onModeChange, discoEnabled, onDiscoEnabledChange, ...props }: Props) {
  return (
    <div className="webagent-composer">
      <div className="webagent-composer-toolbar">
        <label className="form-label" htmlFor="thinking-mode-select">
          Thinking mode
        </label>
        <select
          id="thinking-mode-select"
          className="form-select form-select-sm"
          value={mode}
          onChange={(e) => onModeChange(e.target.value)}
        >
          <option value="fast">Fast</option>
          <option value="deliberate">Deliberate</option>
          <option value="deep">Deep</option>
        </select>
        <div className="form-check form-switch ms-auto">
          <input
            id="disco-dialogue-toggle"
            className="form-check-input"
            type="checkbox"
            checked={discoEnabled}
            onChange={(e) => onDiscoEnabledChange(e.target.checked)}
          />
          <label className="form-check-label" htmlFor="disco-dialogue-toggle">
            Disco dialogue
          </label>
        </div>
      </div>
      <DefaultComposer {...props} />
    </div>
  );
}
