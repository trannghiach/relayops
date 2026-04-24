"use client";

import { api } from "@/lib/api";
import { Card, Col, Row, Typography } from "antd";
import { useEffect, useState } from "react";
import type { MetricsSummary } from "@/lib/types";

export default function OverviewPage() {
  const [metrics, setMetrics] = useState<MetricsSummary | null>(null);

  useEffect(() => {
    api.getMetrics().then((res) => setMetrics(res.data));
  }, []);

  if (!metrics) {
    return <Typography.Text>Loading...</Typography.Text>;
  }

  return (
    <>
      <Typography.Title level={2}>Overview</Typography.Title>

      <Row gutter={[16, 16]}>
        <Col span={6}>
          <Card title="Total Events">{metrics.total_events}</Card>
        </Col>
        <Col span={6}>
          <Card title="Pending Jobs">{metrics.pending_jobs}</Card>
        </Col>
        <Col span={6}>
          <Card title="Succeeded Jobs">{metrics.succeeded_jobs}</Card>
        </Col>
        <Col span={6}>
          <Card title="Dead Letters">{metrics.dead_lettered_jobs}</Card>
        </Col>
        <Col span={6}>
          <Card title="Total Attempts">{metrics.total_attempts}</Card>
        </Col>
      </Row>
    </>
  );
}