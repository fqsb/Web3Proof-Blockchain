import { Button } from "antd";

interface EmptyStateProps {
  icon?: string;
  title: string;
  description?: string;
  actionLabel?: string;
  onAction?: () => void;
}

export default function EmptyState({
  icon = "ID",
  title,
  description,
  actionLabel,
  onAction,
}: EmptyStateProps) {
  return (
    <div className="empty-state">
      <div className="empty-state__icon">{icon}</div>
      <h2 className="empty-state__title">{title}</h2>
      {description && <p className="empty-state__desc">{description}</p>}
      {actionLabel && onAction && (
        <Button type="primary" size="large" onClick={onAction}>
          {actionLabel}
        </Button>
      )}
    </div>
  );
}
