/**
 * 从标题中提取主题信息
 * @param {string} title - 原始标题
 * @returns {{ theme: string, title: string }}
 */
export function getTitleTheme(title) {
    const themePattern = /^\[([iswfISWF])\]\s*(.*)/;
    const match = title.match(themePattern);

    if (!match) {
        return { theme: "", title: title };
    }

    const themeChar = match[1] || '';
    const actualTitle = match[2];
    const theme = themeChar.toLowerCase();
    title = theme ? actualTitle : title;
    return { theme, title };
}

/**
 * TitleInfo 类
 */
export class TitleInfo {
    constructor(theme = "", title = "", user = "", face = "", channel = "") {
        this.theme = theme;
        this.title = title;
        this.user = user;
        this.face = face;
        this.channel = channel;
    }
}

/**
 * 从输入字符串中提取标题信息
 * @param {string} input - 原始输入字符串
 * @returns {TitleInfo}
 */
export function calcTitleInfo(input) {
    if (!input || input.trim() === "") {
        return new TitleInfo();
    }

    const channelPattern = /\[~([^\]]+)\]/;
    const groupPattern = /\[#([^!\]\n]+)(?:!([^\]\n]+))?\]/;

    let title = input.trim();
    let user = "";
    let face = "";
    let channel = "";

    // 提取主题
    const titleTheme = getTitleTheme(title);
    const theme = titleTheme.theme;
    title = titleTheme.title;

    // 提取频道名字
    const channelMatch = title.match(channelPattern);
    if (channelMatch) {
        channel = channelMatch[1]?.trim() || "";
        title = title.replace(channelMatch[0], "").trim();
    }

    // 提取分组信息（用户和头像表情）
    const groupMatch = title.match(groupPattern);
    if (groupMatch) {
        user = groupMatch[1]?.trim() || "";
        face = groupMatch[2]?.trim() || "";
        title = title.replace(groupMatch[0], "").trim();
    }

    // 清理标题中可能剩余的空格
    title = title.trim();

    return new TitleInfo(theme, title, user, face, channel);
}