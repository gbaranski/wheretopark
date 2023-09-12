//
//  WebViewPlugin.swift
//  App
//
//  Created by Grzegorz Barański on 12/09/2023.
//

import Capacitor

@objc(WebViewPlugin)
public class WebViewPlugin: CAPPlugin {
    @objc override public func load() {
        // Called when the plugin is first constructed in the bridge
        self.bridge?.webView?.scrollView.bounces = true
    }
}
