//
//  BottomSheet.swift
//  wheretopark
//
//  Created by Grzegorz Barański on 15/09/2022.
//

import Foundation
import UIKit
import SwiftUI
import Combine

extension UISheetPresentationController.Detent.Identifier {
    public static let small: UISheetPresentationController.Detent.Identifier = UISheetPresentationController.Detent.Identifier(rawValue: "small")
    
    public static let compact: UISheetPresentationController.Detent.Identifier = UISheetPresentationController.Detent.Identifier(rawValue: "compact")
}

extension UISheetPresentationController.Detent {
    class func small() -> UISheetPresentationController.Detent {
        return .custom(identifier: .small) { context in
            return 80
        }
    }
    
    class func compact() -> UISheetPresentationController.Detent {
        return .custom(identifier: .compact) { context in
            return 300
        }
    }
}

class BottomSheetViewController<Content: View>: UIViewController, UISheetPresentationControllerDelegate {
    @Binding private var isPresented: Bool

    @Binding private var selectedDetentIdentifier: UISheetPresentationController.Detent.Identifier?

    private let contentView: UIHostingController<Content>

    init(
        isPresented: Binding<Bool>,
        selectedDetentIdentifier: Binding<UISheetPresentationController.Detent.Identifier?> = Binding.constant(nil),
        isModalInPresentation: Bool = false,
        content: Content
    ) {
        _isPresented = isPresented

        self._selectedDetentIdentifier = selectedDetentIdentifier
        
        self.contentView = UIHostingController(rootView: content)

        super.init(nibName: nil, bundle: nil)
        self.isModalInPresentation = isModalInPresentation
    }

    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    override func viewDidLoad() {
        super.viewDidLoad()

        addChild(contentView)
        view.addSubview(contentView.view)

        contentView.view.translatesAutoresizingMaskIntoConstraints = false

        NSLayoutConstraint.activate([
            contentView.view.topAnchor.constraint(equalTo: view.topAnchor),
            contentView.view.bottomAnchor.constraint(equalTo: view.bottomAnchor),
            contentView.view.leadingAnchor.constraint(equalTo: view.leadingAnchor),
            contentView.view.trailingAnchor.constraint(equalTo: view.trailingAnchor)
        ])

        if let presentationController = presentationController as? UISheetPresentationController {
            presentationController.detents = [.small(), .compact(), .large()]
            presentationController.largestUndimmedDetentIdentifier = .large
            presentationController.prefersGrabberVisible = true
            presentationController.prefersScrollingExpandsWhenScrolledToEdge = true
            presentationController.prefersEdgeAttachedInCompactHeight = false
            presentationController.selectedDetentIdentifier = selectedDetentIdentifier
            presentationController.widthFollowsPreferredContentSizeWhenEdgeAttached = false
            presentationController.delegate = self
        }
    }

    override func viewDidDisappear(_ animated: Bool) {
        super.viewDidDisappear(animated)

        isPresented = false
    }
    
    func updateSelectedDetentIdentifier(_ selectedDetentIdentifier: UISheetPresentationController.Detent.Identifier?) {
        self.sheetPresentationController?.animateChanges {
            self.sheetPresentationController?.selectedDetentIdentifier = selectedDetentIdentifier
        }
    }
    
    func sheetPresentationControllerDidChangeSelectedDetentIdentifier(_ sheetPresentationController: UISheetPresentationController) {
        self.selectedDetentIdentifier = sheetPresentationController.selectedDetentIdentifier
    }
}


struct BottomSheet<T: Any, ContentView: View>: ViewModifier {
    @Binding private var isPresented: Bool
    
    @Binding private var selectedDetentIdentifier: UISheetPresentationController.Detent.Identifier?
    private let isModalInPresentation: Bool
    private var onDismiss: (() -> Void)?
    private let contentView: () -> ContentView
    
    @State private var bottomSheetViewController: BottomSheetViewController<ContentView>?

    init(
        isPresented: Binding<Bool>,
        selectedDetentIdentifier: Binding<UISheetPresentationController.Detent.Identifier?> = Binding.constant(nil),
        isModalInPresentation: Bool = false,
        onDismiss: (() -> Void)? = nil,
        @ViewBuilder contentView: @escaping () -> ContentView
    ) {
        _isPresented = isPresented
        self._selectedDetentIdentifier = selectedDetentIdentifier
        self.contentView = contentView
        self.onDismiss = onDismiss
        self.isModalInPresentation = isModalInPresentation
    }
    
    func body(content: Content) -> some View {
        content
            .onChange(of: isPresented, perform: updatePresentation)
            .onChange(of: selectedDetentIdentifier, perform: updateSelectedDetentIdentifier)
    }

    private func updatePresentation(_ isPresented: Bool) {
        guard let windowScene = UIApplication.shared.connectedScenes.first(where: {
            $0.activationState == .foregroundActive
        }) as? UIWindowScene else { return }

        
        guard let root = windowScene.keyWindow?.rootViewController else { return }
        var controllerToPresentFrom = root
        while let presented = controllerToPresentFrom.presentedViewController {
            controllerToPresentFrom = presented
        }

        if isPresented {
            bottomSheetViewController = BottomSheetViewController(
                isPresented: $isPresented,
                selectedDetentIdentifier: $selectedDetentIdentifier,
                isModalInPresentation: isModalInPresentation,
                content: contentView()
            )

            controllerToPresentFrom.present(bottomSheetViewController!, animated: true)

        } else {
            bottomSheetViewController?.dismiss(animated: true, completion: onDismiss)
        }
    }
    
    private func updateSelectedDetentIdentifier(_ selectedDetentIdentifier: UISheetPresentationController.Detent.Identifier?) {
        bottomSheetViewController?.updateSelectedDetentIdentifier(selectedDetentIdentifier)
    }
}

extension View {
    public func bottomSheet<ContentView: View>(
        isPresented: Binding<Bool>,
        selectedDetentIdentifier: Binding<UISheetPresentationController.Detent.Identifier?> = Binding.constant(nil),
        isModalInPresentation: Bool = false,
        onDismiss: (() -> Void)? = nil,
        @ViewBuilder contentView: @escaping () -> ContentView
    ) -> some View {
        self.modifier(
            BottomSheet<Any, ContentView>(
                isPresented: isPresented,
                selectedDetentIdentifier: selectedDetentIdentifier,
                isModalInPresentation: isModalInPresentation,
                onDismiss: onDismiss,
                contentView: contentView
            )
        )
    }
}
