//
//  ViewController.swift
//  Networth
//
//  Created by Kien Pham on 8/16/18.
//  Copyright © 2018 Networth. All rights reserved.
//

import Cocoa
import AppKit
import WebKit

class ViewController: NSViewController {
    
    @IBOutlet weak var webView: WKWebView!
    @IBOutlet weak var address: NSTextField!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        let addressStr:String = "https://apple.com"
        let url:URL = URL(string: addressStr)!
        let urlRequest:URLRequest = URLRequest(url: url)
        webView.load(urlRequest)
        
//        address.value = addressStr
    }
    
//    override func viewDidAppear() {
//        super.viewDidAppear()
//
////        let url:URL = URL(string: "https://apple.com")!
////
////        let urlRequest:URLRequest = URLRequest(url: url)
////        webView.load(urlRequest)
//    }
    
    override var representedObject: Any? {
        didSet {
        // Update the view, if already loaded.
        }
    }

//    func text

}

