"use client";

import {
    DashboardOutlined,
    DatabaseOutlined,
    WarningOutlined,
    ProfileOutlined,
} from "@ant-design/icons";
import { Layout, Menu, Typography } from "antd";
import Link from "next/link";
import { usePathname } from "next/navigation";
import type { ReactNode } from "react";

const { Header, Sider, Content } = Layout;

export function AppShell({ children }: { children: ReactNode }) {
    const pathname = usePathname();

    return (
        <Layout style={{ minHeight: "100vh" }}>
            <Sider width={240}>
                <div style={{ padding: 20 }}>
                    <Typography.Title level={4} style={{ color: "white", margin: 0 }}>
                        RelayOps
                    </Typography.Title>
                </div>

                <Menu
                    theme="dark"
                    mode="inline"
                    selectedKeys={[pathname]}
                    items={[
                        {
                            key: "/",
                            icon: <DashboardOutlined />,
                            label: <Link href="/">Overview</Link>,
                        },
                        {
                            key: "/events",
                            icon: <ProfileOutlined />,
                            label: <Link href="/events">Events</Link>,
                        },
                        {
                            key: "/jobs",
                            icon: <DatabaseOutlined />,
                            label: <Link href="/jobs">Jobs</Link>,
                        },
                        {
                            key: "/dead-letters",
                            icon: <WarningOutlined />,
                            label: <Link href="/dead-letters">Dead Letters</Link>,
                        },
                    ]}
                />
            </Sider>

            <Layout>
                <Header style={{ background: "white" }}>
                    <Typography.Text strong>
                        Event-driven notification & workflow platform
                    </Typography.Text>
                </Header>

                <Content style={{ padding: 24 }}>{children}</Content>
            </Layout>
        </Layout>
    );
}