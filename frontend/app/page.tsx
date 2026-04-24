"use client";

import { api } from "@/lib/api";
import type { EventItem, JobItem, MetricsSummary } from "@/lib/types";
import { Card, Col, Row, Space, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import Link from "next/link";
import { useEffect, useState } from "react";

function statusColor(status: string) {
  switch (status) {
    case "succeeded":
    case "planned":
      return "green";
    case "pending":
      return "blue";
    case "processing":
      return "orange";
    case "dead_lettered":
    case "failed":
      return "red";
    default:
      return "default";
  }
}

export default function OverviewPage() {
  const [metrics, setMetrics] = useState<MetricsSummary | null>(null);
  const [events, setEvents] = useState<EventItem[]>([]);
  const [jobs, setJobs] = useState<JobItem[]>([]);

  useEffect(() => {
    Promise.all([api.getMetrics(), api.getEvents(), api.getJobs()]).then(
      ([metricsRes, eventsRes, jobsRes]) => {
        setMetrics(metricsRes.data);
        setEvents(eventsRes.data.slice(0, 5));
        setJobs(jobsRes.data.slice(0, 5));
      }
    );
  }, []);

  if (!metrics) {
    return <Typography.Text>Loading...</Typography.Text>;
  }

  const eventColumns: ColumnsType<EventItem> = [
    { title: "Type", dataIndex: "event_type" },
    { title: "Source", dataIndex: "source" },
    {
      title: "Status",
      dataIndex: "status",
      render: (status) => <Tag color={statusColor(status)}>{status}</Tag>,
    },
    {
      title: "Created At",
      dataIndex: "created_at",
      render: (v) => new Date(v).toLocaleString(),
    },
  ];

  const jobColumns: ColumnsType<JobItem> = [
    { title: "Type", dataIndex: "job_type" },
    { title: "Channel", dataIndex: "channel" },
    {
      title: "Status",
      dataIndex: "status",
      render: (status) => <Tag color={statusColor(status)}>{status}</Tag>,
    },
    {
      title: "Executions",
      render: (_, row) => `${row.attempts} total`,
    },
    {
      title: "Action",
      render: (_, row) => <Link href={`/jobs/${row.id}`}>View</Link>,
    },
  ];

  return (
    <Space orientation="vertical" size="large" style={{ width: "100%" }}>
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

      <Card title="Recent Jobs">
        <Table
          rowKey="id"
          columns={jobColumns}
          dataSource={jobs}
          pagination={false}
        />
      </Card>

      <Card title="Recent Events">
        <Table
          rowKey="id"
          columns={eventColumns}
          dataSource={events}
          pagination={false}
        />
      </Card>
    </Space>
  );
}